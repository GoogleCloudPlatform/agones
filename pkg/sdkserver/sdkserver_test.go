// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdkserver

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"agones.dev/agones/pkg/apis/stable/v1alpha1"
	"agones.dev/agones/pkg/sdk"
	agtesting "agones.dev/agones/pkg/testing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

func TestSidecarRun(t *testing.T) {
	t.Parallel()

	type expected struct {
		state       v1alpha1.GameServerState
		labels      map[string]string
		annotations map[string]string
		recordings  []string
	}

	fixtures := map[string]struct {
		f        func(*SDKServer, context.Context)
		expected expected
	}{
		"ready": {
			f: func(sc *SDKServer, ctx context.Context) {
				sc.Ready(ctx, &sdk.Empty{}) // nolint: errcheck
			},
			expected: expected{
				state: v1alpha1.GameServerStateRequestReady,
			},
		},
		"shutdown": {
			f: func(sc *SDKServer, ctx context.Context) {
				sc.Shutdown(ctx, &sdk.Empty{}) // nolint: errcheck
			},
			expected: expected{
				state: v1alpha1.GameServerStateShutdown,
			},
		},
		"unhealthy": {
			f: func(sc *SDKServer, ctx context.Context) {
				// we have a 1 second timeout
				time.Sleep(2 * time.Second)
			},
			expected: expected{
				state:      v1alpha1.GameServerStateUnhealthy,
				recordings: []string{string(v1alpha1.GameServerStateUnhealthy)},
			},
		},
		"label": {
			f: func(sc *SDKServer, ctx context.Context) {
				_, err := sc.SetLabel(ctx, &sdk.KeyValue{Key: "foo", Value: "value-foo"})
				assert.Nil(t, err)
				_, err = sc.SetLabel(ctx, &sdk.KeyValue{Key: "bar", Value: "value-bar"})
				assert.Nil(t, err)
			},
			expected: expected{
				labels: map[string]string{
					metadataPrefix + "foo": "value-foo",
					metadataPrefix + "bar": "value-bar"},
			},
		},
		"annotation": {
			f: func(sc *SDKServer, ctx context.Context) {
				_, err := sc.SetAnnotation(ctx, &sdk.KeyValue{Key: "test-1", Value: "annotation-1"})
				assert.Nil(t, err)
				_, err = sc.SetAnnotation(ctx, &sdk.KeyValue{Key: "test-2", Value: "annotation-2"})
				assert.Nil(t, err)
			},
			expected: expected{
				annotations: map[string]string{
					metadataPrefix + "test-1": "annotation-1",
					metadataPrefix + "test-2": "annotation-2"},
			},
		},
		"allocated": {
			f: func(sc *SDKServer, ctx context.Context) {
				_, err := sc.Allocate(ctx, &sdk.Empty{})
				assert.NoError(t, err)
			},
			expected: expected{
				state: v1alpha1.GameServerStateAllocated,
			},
		},
	}

	for k, v := range fixtures {
		t.Run(k, func(t *testing.T) {
			m := agtesting.NewMocks()
			done := make(chan bool)

			m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				gs := v1alpha1.GameServer{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test", Namespace: "default",
					},
					Spec: v1alpha1.GameServerSpec{
						Health: v1alpha1.Health{Disabled: false, FailureThreshold: 1, PeriodSeconds: 1, InitialDelaySeconds: 0},
					},
					Status: v1alpha1.GameServerStatus{
						State: v1alpha1.GameServerStateStarting,
					},
				}
				gs.ApplyDefaults()
				return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{gs}}, nil
			})
			m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				defer close(done)
				ua := action.(k8stesting.UpdateAction)
				gs := ua.GetObject().(*v1alpha1.GameServer)

				if v.expected.state != "" {
					assert.Equal(t, v.expected.state, gs.Status.State)
				}

				for label, value := range v.expected.labels {
					assert.Equal(t, value, gs.ObjectMeta.Labels[label])
				}
				for ann, value := range v.expected.annotations {
					assert.Equal(t, value, gs.ObjectMeta.Annotations[ann])
				}

				return true, gs, nil
			})

			sc, err := NewSDKServer("test", "default", m.KubeClient, m.AgonesClient)
			stop := make(chan struct{})
			defer close(stop)
			sc.informerFactory.Start(stop)
			assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))

			assert.Nil(t, err)
			sc.recorder = m.FakeRecorder

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
				err := sc.Run(ctx.Done())
				assert.Nil(t, err)
				wg.Done()
			}()
			v.f(sc, ctx)

			select {
			case <-done:
			case <-time.After(10 * time.Second):
				assert.Fail(t, "Timeout on Run")
			}

			logrus.Info("attempting to find event recording")
			for _, str := range v.expected.recordings {
				agtesting.AssertEventContains(t, m.FakeRecorder.Events, str)
			}

			cancel()
			wg.Wait()
		})
	}
}

func TestSDKServerSyncGameServer(t *testing.T) {
	t.Parallel()

	type expected struct {
		state       v1alpha1.GameServerState
		labels      map[string]string
		annotations map[string]string
	}

	type scData struct {
		gsState       v1alpha1.GameServerState
		gsLabels      map[string]string
		gsAnnotations map[string]string
	}

	fixtures := map[string]struct {
		expected expected
		key      string
		scData   scData
	}{
		"ready": {
			key: string(updateState),
			scData: scData{
				gsState: v1alpha1.GameServerStateReady,
			},
			expected: expected{
				state: v1alpha1.GameServerStateReady,
			},
		},
		"label": {
			key: string(updateLabel),
			scData: scData{
				gsLabels: map[string]string{"foo": "bar"},
			},
			expected: expected{
				labels: map[string]string{metadataPrefix + "foo": "bar"},
			},
		},
		"annotation": {
			key: string(updateAnnotation),
			scData: scData{
				gsAnnotations: map[string]string{"test": "annotation"},
			},
			expected: expected{
				annotations: map[string]string{metadataPrefix + "test": "annotation"},
			},
		},
	}

	for k, v := range fixtures {
		t.Run(k, func(t *testing.T) {
			m := agtesting.NewMocks()
			sc, err := defaultSidecar(m)
			assert.Nil(t, err)

			sc.gsState = v.scData.gsState
			sc.gsLabels = v.scData.gsLabels
			sc.gsAnnotations = v.scData.gsAnnotations

			updated := false

			m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				gs := v1alpha1.GameServer{ObjectMeta: metav1.ObjectMeta{
					UID:  "1234",
					Name: sc.gameServerName, Namespace: sc.namespace,
					Labels: map[string]string{}, Annotations: map[string]string{}},
				}
				return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{gs}}, nil
			})
			m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				updated = true
				ua := action.(k8stesting.UpdateAction)
				gs := ua.GetObject().(*v1alpha1.GameServer)

				if v.expected.state != "" {
					assert.Equal(t, v.expected.state, gs.Status.State)
				}
				for label, value := range v.expected.labels {
					assert.Equal(t, value, gs.ObjectMeta.Labels[label])
				}
				for ann, value := range v.expected.annotations {
					assert.Equal(t, value, gs.ObjectMeta.Annotations[ann])
				}

				return true, gs, nil
			})

			stop := make(chan struct{})
			defer close(stop)
			sc.informerFactory.Start(stop)
			assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))

			err = sc.syncGameServer(v.key)
			assert.Nil(t, err)
			assert.True(t, updated, "should have updated")

		})
	}
}

func TestSidecarUpdateState(t *testing.T) {
	t.Parallel()

	fixtures := map[string]struct {
		f func(gs *v1alpha1.GameServer)
	}{
		"unhealthy": {
			f: func(gs *v1alpha1.GameServer) {
				gs.Status.State = v1alpha1.GameServerStateUnhealthy
			},
		},
		"shutdown": {
			f: func(gs *v1alpha1.GameServer) {
				gs.Status.State = v1alpha1.GameServerStateShutdown
			},
		},
		"deleted": {
			f: func(gs *v1alpha1.GameServer) {
				now := metav1.Now()
				gs.ObjectMeta.DeletionTimestamp = &now
			},
		},
	}

	for k, v := range fixtures {
		t.Run(k, func(t *testing.T) {
			m := agtesting.NewMocks()
			sc, err := defaultSidecar(m)
			assert.Nil(t, err)
			sc.gsState = v1alpha1.GameServerStateReady

			updated := false

			m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				gs := v1alpha1.GameServer{
					ObjectMeta: metav1.ObjectMeta{Name: sc.gameServerName, Namespace: sc.namespace},
					Status:     v1alpha1.GameServerStatus{},
				}

				// apply mutation
				v.f(&gs)

				return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{gs}}, nil
			})
			m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
				updated = true
				return true, nil, nil
			})

			stop := make(chan struct{})
			defer close(stop)
			sc.informerFactory.Start(stop)
			assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))

			err = sc.updateState()
			assert.Nil(t, err)
			assert.False(t, updated)
		})
	}
}

func TestSidecarHealthLastUpdated(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	m := agtesting.NewMocks()

	sc, err := defaultSidecar(m)
	assert.Nil(t, err)
	sc.health = v1alpha1.Health{Disabled: false}
	fc := clock.NewFakeClock(now)
	sc.clock = fc

	stream := newEmptyMockStream()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := sc.Health(stream)
		assert.Nil(t, err)
		wg.Done()
	}()

	// Test once with a single message
	fc.Step(3 * time.Second)
	stream.msgs <- &sdk.Empty{}

	err = waitForMessage(sc)
	assert.Nil(t, err)
	sc.healthMutex.RLock()
	assert.Equal(t, sc.clock.Now().UTC().String(), sc.healthLastUpdated.String())
	sc.healthMutex.RUnlock()

	// Test again, since the value has been set, that it is re-set
	fc.Step(3 * time.Second)
	stream.msgs <- &sdk.Empty{}
	err = waitForMessage(sc)
	assert.Nil(t, err)
	sc.healthMutex.RLock()
	assert.Equal(t, sc.clock.Now().UTC().String(), sc.healthLastUpdated.String())
	sc.healthMutex.RUnlock()

	// make sure closing doesn't change the time
	fc.Step(3 * time.Second)
	close(stream.msgs)
	assert.NotEqual(t, sc.clock.Now().UTC().String(), sc.healthLastUpdated.String())

	wg.Wait()
}

func TestSidecarHealthy(t *testing.T) {
	t.Parallel()

	m := agtesting.NewMocks()
	sc, err := defaultSidecar(m)
	// manually set the values
	sc.health = v1alpha1.Health{FailureThreshold: 1}
	sc.healthTimeout = 5 * time.Second
	sc.initHealthLastUpdated(0 * time.Second)

	assert.Nil(t, err)

	now := time.Now().UTC()
	fc := clock.NewFakeClock(now)
	sc.clock = fc

	stream := newEmptyMockStream()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := sc.Health(stream)
		assert.Nil(t, err)
		wg.Done()
	}()

	fixtures := map[string]struct {
		timeAdd         time.Duration
		disabled        bool
		expectedHealthy bool
	}{
		"disabled, under timeout": {disabled: true, timeAdd: time.Second, expectedHealthy: true},
		"disabled, over timeout":  {disabled: true, timeAdd: 15 * time.Second, expectedHealthy: true},
		"enabled, under timeout":  {disabled: false, timeAdd: time.Second, expectedHealthy: true},
		"enabled, over timeout":   {disabled: false, timeAdd: 15 * time.Second, expectedHealthy: false},
	}

	for k, v := range fixtures {
		t.Run(k, func(t *testing.T) {
			logrus.WithField("test", k).Infof("Test Running")
			sc.health.Disabled = v.disabled
			fc.SetTime(time.Now().UTC())
			stream.msgs <- &sdk.Empty{}
			err = waitForMessage(sc)
			assert.Nil(t, err)

			fc.Step(v.timeAdd)
			sc.checkHealth()
			assert.Equal(t, v.expectedHealthy, sc.healthy())
		})
	}

	t.Run("initial delay", func(t *testing.T) {
		sc.health.Disabled = false
		fc.SetTime(time.Now().UTC())
		sc.initHealthLastUpdated(0)
		sc.healthFailureCount = 0
		sc.checkHealth()
		assert.True(t, sc.healthy())

		sc.initHealthLastUpdated(10 * time.Second)
		sc.checkHealth()
		assert.True(t, sc.healthy())
		fc.Step(9 * time.Second)
		sc.checkHealth()
		assert.True(t, sc.healthy())

		fc.Step(10 * time.Second)
		sc.checkHealth()
		assert.False(t, sc.healthy())
	})

	t.Run("health failure threshold", func(t *testing.T) {
		sc.health.Disabled = false
		sc.health.FailureThreshold = 3
		fc.SetTime(time.Now().UTC())
		sc.initHealthLastUpdated(0)
		sc.healthFailureCount = 0

		sc.checkHealth()
		assert.True(t, sc.healthy())
		assert.Equal(t, int32(0), sc.healthFailureCount)

		fc.Step(10 * time.Second)
		sc.checkHealth()
		assert.True(t, sc.healthy())
		sc.checkHealth()
		assert.True(t, sc.healthy())
		sc.checkHealth()
		assert.False(t, sc.healthy())

		stream.msgs <- &sdk.Empty{}
		err = waitForMessage(sc)
		assert.Nil(t, err)
		fc.Step(10 * time.Second)
		assert.True(t, sc.healthy())
	})

	close(stream.msgs)
	wg.Wait()
}

func TestSidecarHTTPHealthCheck(t *testing.T) {
	m := agtesting.NewMocks()
	sc, err := NewSDKServer("test", "default", m.KubeClient, m.AgonesClient)
	assert.Nil(t, err)
	now := time.Now().Add(time.Hour).UTC()
	fc := clock.NewFakeClock(now)
	// now we control time - so slow machines won't fail anymore
	sc.clock = fc

	m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
		gs := v1alpha1.GameServer{
			ObjectMeta: metav1.ObjectMeta{Name: sc.gameServerName, Namespace: sc.namespace},
			Spec: v1alpha1.GameServerSpec{
				Health: v1alpha1.Health{Disabled: false, FailureThreshold: 1, PeriodSeconds: 1, InitialDelaySeconds: 0},
			},
		}

		return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{gs}}, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)

	step := 2 * time.Second

	go func() {
		err := sc.Run(ctx.Done())
		assert.Nil(t, err)
		// gate
		assert.Equal(t, 1*time.Second, sc.healthTimeout)
		wg.Done()
	}()

	testHTTPHealth(t, "http://localhost:8080/healthz", "ok", http.StatusOK)
	testHTTPHealth(t, "http://localhost:8080/gshealthz", "ok", http.StatusOK)

	assert.Equal(t, now, sc.healthLastUpdated)

	fc.Step(step)
	time.Sleep(step)
	assert.False(t, sc.healthy())
	testHTTPHealth(t, "http://localhost:8080/gshealthz", "", http.StatusInternalServerError)
	cancel()
	wg.Wait() // wait for go routine test results.
}

func TestSDKServerGetGameServer(t *testing.T) {
	t.Parallel()

	fixture := &v1alpha1.GameServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.GameServerStatus{
			State: v1alpha1.GameServerStateReady,
		},
	}

	m := agtesting.NewMocks()
	m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{*fixture}}, nil
	})

	stop := make(chan struct{})
	defer close(stop)

	sc, err := defaultSidecar(m)
	assert.Nil(t, err)

	sc.informerFactory.Start(stop)
	assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))

	result, err := sc.GetGameServer(context.Background(), &sdk.Empty{})
	assert.Nil(t, err)
	assert.Equal(t, fixture.ObjectMeta.Name, result.ObjectMeta.Name)
	assert.Equal(t, fixture.ObjectMeta.Namespace, result.ObjectMeta.Namespace)
	assert.Equal(t, string(fixture.Status.State), result.Status.State)
}

func TestSDKServerWatchGameServer(t *testing.T) {
	t.Parallel()
	m := agtesting.NewMocks()
	sc, err := defaultSidecar(m)
	assert.Nil(t, err)
	assert.Empty(t, sc.connectedStreams)

	stream := newGameServerMockStream()
	asyncWatchGameServer(t, sc, stream)
	assert.Nil(t, waitConnectedStreamCount(sc, 1))
	assert.Equal(t, stream, sc.connectedStreams[0])

	stream = newGameServerMockStream()
	asyncWatchGameServer(t, sc, stream)
	assert.Nil(t, waitConnectedStreamCount(sc, 2))
	assert.Len(t, sc.connectedStreams, 2)
	assert.Equal(t, stream, sc.connectedStreams[1])
}

func TestSDKServerSendGameServerUpdate(t *testing.T) {
	t.Parallel()
	m := agtesting.NewMocks()
	sc, err := defaultSidecar(m)
	assert.Nil(t, err)
	assert.Empty(t, sc.connectedStreams)

	stop := make(chan struct{})
	defer close(stop)
	sc.stop = stop

	stream := newGameServerMockStream()
	asyncWatchGameServer(t, sc, stream)
	assert.Nil(t, waitConnectedStreamCount(sc, 1))

	fixture := &v1alpha1.GameServer{ObjectMeta: metav1.ObjectMeta{Name: "test-server"}}

	sc.sendGameServerUpdate(fixture)

	var sdkGS *sdk.GameServer
	select {
	case sdkGS = <-stream.msgs:
	case <-time.After(3 * time.Second):
		assert.Fail(t, "Event stream should not have timed out")
	}

	assert.Equal(t, fixture.ObjectMeta.Name, sdkGS.ObjectMeta.Name)
}

func TestSDKServerUpdateEventHandler(t *testing.T) {
	t.Parallel()
	m := agtesting.NewMocks()

	fakeWatch := watch.NewFake()
	m.AgonesClient.AddWatchReactor("gameservers", k8stesting.DefaultWatchReactor(fakeWatch, nil))

	stop := make(chan struct{})
	defer close(stop)

	sc, err := defaultSidecar(m)
	assert.Nil(t, err)

	sc.informerFactory.Start(stop)
	assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))

	stream := newGameServerMockStream()
	asyncWatchGameServer(t, sc, stream)
	assert.Nil(t, waitConnectedStreamCount(sc, 1))

	fixture := &v1alpha1.GameServer{ObjectMeta: metav1.ObjectMeta{Name: "test-server", Namespace: "default"},
		Spec: v1alpha1.GameServerSpec{},
	}

	// need to add it before it can be modified
	fakeWatch.Add(fixture.DeepCopy())
	fakeWatch.Modify(fixture.DeepCopy())

	var sdkGS *sdk.GameServer
	select {
	case sdkGS = <-stream.msgs:
	case <-time.After(3 * time.Second):
		assert.Fail(t, "Event stream should not have timed out")
	}

	assert.NotNil(t, sdkGS)
	assert.Equal(t, fixture.ObjectMeta.Name, sdkGS.ObjectMeta.Name)
}

func TestSDKServerAllocate(t *testing.T) {
	t.Parallel()

	defaultGs := &v1alpha1.GameServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.GameServerStatus{
			State: v1alpha1.GameServerStateReady,
		},
	}

	setup := func(t *testing.T, fixture *v1alpha1.GameServer) (agtesting.Mocks, *SDKServer, chan struct{}) {
		m := agtesting.NewMocks()
		m.AgonesClient.AddReactor("list", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, &v1alpha1.GameServerList{Items: []v1alpha1.GameServer{*fixture}}, nil
		})
		stop := make(chan struct{})
		sc, err := defaultSidecar(m)
		assert.Nil(t, err)
		sc.informerFactory.Start(stop)
		assert.True(t, cache.WaitForCacheSync(stop, sc.gameServerSynced))
		return m, sc, stop
	}

	t.Run("works perfectly", func(t *testing.T) {
		m, sc, stop := setup(t, defaultGs.DeepCopy())
		defer close(stop)

		done := make(chan struct{}, 10)
		m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
			ua := action.(k8stesting.UpdateAction)
			gs := ua.GetObject().(*v1alpha1.GameServer)
			assert.Equal(t, v1alpha1.GameServerStateAllocated, gs.Status.State)

			defer func() {
				done <- struct{}{}
			}()

			return true, gs, nil
		})

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.Fail(t, "should be updated")
		}
	})

	t.Run("contention on first build", func(t *testing.T) {
		m, sc, stop := setup(t, defaultGs.DeepCopy())
		defer close(stop)

		done := make(chan struct{}, 10)
		count := 0

		m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
			ua := action.(k8stesting.UpdateAction)
			defer func() {
				done <- struct{}{}
			}()

			count++
			if count == 1 {
				return true, nil, k8serrors.NewConflict(schema.ParseGroupResource(""), "gameservers", errors.New("contention"))
			}

			gs := ua.GetObject().(*v1alpha1.GameServer)
			assert.Equal(t, v1alpha1.GameServerStateAllocated, gs.Status.State)
			return true, gs, nil
		})

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.Fail(t, "have contention")
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.Fail(t, "should be updated")
		}
	})

	t.Run("timeout", func(t *testing.T) {
		m, sc, stop := setup(t, defaultGs.DeepCopy())
		defer close(stop)

		now := time.Now()
		fc := clock.NewFakeClock(now)
		sc.clock = fc

		done := make(chan struct{}, 10)

		m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
			defer func() {
				done <- struct{}{}
			}()

			// move past timeout
			fc.Step(40 * time.Second)
			return true, nil, k8serrors.NewConflict(schema.ParseGroupResource(""), "gameservers", errors.New("contention"))
		})

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.EqualError(t, err, "Allocation request timed out")

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.Fail(t, "have contention")
		}
	})

	t.Run("unhealthy gameserver", func(t *testing.T) {
		gs := defaultGs.DeepCopy()
		gs.Status.State = v1alpha1.GameServerStateUnhealthy
		_, sc, stop := setup(t, gs.DeepCopy())
		defer close(stop)

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.EqualError(t, err, "cannot Allocate an Unhealthy GameServer")
	})

	t.Run("Shutdown gameserver", func(t *testing.T) {
		gs := defaultGs.DeepCopy()
		gs.Status.State = v1alpha1.GameServerStateShutdown
		_, sc, stop := setup(t, gs.DeepCopy())
		defer close(stop)

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.EqualError(t, err, "cannot Allocate a Shutdown GameServer")
	})

	t.Run("Deleting gameserver", func(t *testing.T) {
		gs := defaultGs.DeepCopy()
		now := metav1.Now()
		gs.ObjectMeta.DeletionTimestamp = &now
		_, sc, stop := setup(t, gs.DeepCopy())
		defer close(stop)

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.NoError(t, err)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		m, sc, stop := setup(t, defaultGs.DeepCopy())
		defer close(stop)
		done := make(chan struct{}, 10)

		m.AgonesClient.AddReactor("update", "gameservers", func(action k8stesting.Action) (bool, runtime.Object, error) {
			defer func() {
				done <- struct{}{}
			}()

			return true, nil, errors.New("a bad thing happened")
		})

		_, err := sc.Allocate(context.Background(), &sdk.Empty{})
		assert.Error(t, err)
	})
}

func defaultSidecar(m agtesting.Mocks) (*SDKServer, error) {
	return NewSDKServer("test", "default", m.KubeClient, m.AgonesClient)
}

func waitForMessage(sc *SDKServer) error {
	return wait.PollImmediate(time.Second, 5*time.Second, func() (done bool, err error) {
		sc.healthMutex.RLock()
		defer sc.healthMutex.RUnlock()
		return sc.clock.Now().UTC() == sc.healthLastUpdated, nil
	})
}

func waitConnectedStreamCount(sc *SDKServer, count int) error {
	return wait.PollImmediate(1*time.Second, 10*time.Second, func() (bool, error) {
		sc.streamMutex.RLock()
		defer sc.streamMutex.RUnlock()
		return len(sc.connectedStreams) == count, nil
	})
}

func asyncWatchGameServer(t *testing.T, sc *SDKServer, stream sdk.SDK_WatchGameServerServer) {
	go func() {
		err := sc.WatchGameServer(&sdk.Empty{}, stream)
		assert.Nil(t, err)
	}()
}
