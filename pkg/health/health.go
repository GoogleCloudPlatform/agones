// Copyright 2018 Google Inc. All Rights Reserved.
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

package health

import (
	"time"

	"agones.dev/agones/pkg/apis/stable"
	"agones.dev/agones/pkg/apis/stable/v1alpha1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	getterv1alpha1 "agones.dev/agones/pkg/client/clientset/versioned/typed/stable/v1alpha1"
	"agones.dev/agones/pkg/client/informers/externalversions"
	listerv1alpha1 "agones.dev/agones/pkg/client/listers/stable/v1alpha1"
	"agones.dev/agones/pkg/util/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisterv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

// DefaultMonitor watches Pods, and applies
// an Unhealthy state if the GameServer main container exits when in
// a Ready state
type DefaultMonitor struct {
	podSynced        cache.InformerSynced
	podLister        corelisterv1.PodLister
	gameServerGetter getterv1alpha1.GameServersGetter
	gameServerLister listerv1alpha1.GameServerLister
	queue            workqueue.RateLimitingInterface
	recorder         record.EventRecorder
}

// Monitor describes the health interface.
type Monitor interface {
	Run(stop <-chan struct{})
}

// NewMonitor is the default initializer for a DefaultMonitor
func NewMonitor(kubeClient kubernetes.Interface, agonesClient versioned.Interface, kubeInformerFactory informers.SharedInformerFactory,
	agonesInformerFactory externalversions.SharedInformerFactory) *DefaultMonitor {

	podInformer := kubeInformerFactory.Core().V1().Pods().Informer()
	m := &DefaultMonitor{
		podSynced:        podInformer.HasSynced,
		podLister:        kubeInformerFactory.Core().V1().Pods().Lister(),
		gameServerGetter: agonesClient.StableV1alpha1(),
		gameServerLister: agonesInformerFactory.Stable().V1alpha1().GameServers().Lister(),
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), stable.GroupName+".HealthMonitor"),
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	m.recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "health-monitor"})

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			if owner := metav1.GetControllerOf(pod); owner != nil && owner.Kind == "GameServer" {
				if v1alpha1.GameServerRolePodSelector.Matches(labels.Set(pod.Labels)) && m.failedContainer(pod) {
					key := pod.ObjectMeta.Namespace + "/" + owner.Name
					logrus.WithField("key", key).Info("GameServer container has terminated")
					m.enqueueGameServer(key)
				}
			}
		},
	})
	return m
}

// failedContainer checks each container, and determines if there was a failed
// container
func (m *DefaultMonitor) failedContainer(pod *corev1.Pod) bool {
	container := pod.Annotations[v1alpha1.GameServerContainerAnnotation]
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name == container && cs.State.Terminated != nil {
			return true
		}
	}
	return false

}

// enqueue puts the name of the GameServer into the queue
func (hc *DefaultMonitor) enqueueGameServer(key string) {
	hc.queue.AddRateLimited(key)
}

// Run processes the rate limited queue.
// Will block until stop is closed
func (hc *DefaultMonitor) Run(stop <-chan struct{}) {
	defer hc.queue.ShutDown()

	logrus.Info("Starting health worker")
	go wait.Until(hc.runWorker, time.Second, stop)
	<-stop
	logrus.Info("Shut down health worker")
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (hc *DefaultMonitor) runWorker() {
	for hc.processNextWorkItem() {
	}
}

func (hc *DefaultMonitor) processNextWorkItem() bool {
	obj, quit := hc.queue.Get()
	if quit {
		return false
	}
	defer hc.queue.Done(obj)

	logrus.WithField("obj", obj).Info("Processing obj")

	var key string
	var ok bool
	if key, ok = obj.(string); !ok {
		runtime.HandleError(logrus.WithField("obj", obj), errors.Errorf("expected string in queue, but got %T", obj))
		// this is a bad entry, we don't want to reprocess
		hc.queue.Forget(obj)
		return true
	}

	if err := hc.syncGameServer(key); err != nil {
		// we don't forget here, because we want this to be retried via the queue
		runtime.HandleError(logrus.WithField("obj", obj), err)
		hc.queue.AddRateLimited(obj)
		return true
	}

	hc.queue.Forget(obj)
	return true
}

// syncGameServer sets the GameSerer to Unhealthy, if its state is Ready
func (hc *DefaultMonitor) syncGameServer(key string) error {
	logrus.WithField("key", key).Info("Synchronising")

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		// don't return an error, as we don't want this retried
		runtime.HandleError(logrus.WithField("key", key), errors.Wrapf(err, "invalid resource key"))
		return nil
	}

	gs, err := hc.gameServerLister.GameServers(namespace).Get(name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			logrus.WithField("key", key).Info("GameServer is no longer available for syncing")
			return nil
		}
		return errors.Wrapf(err, "error retrieving GameServer %s from namespace %s", name, namespace)
	}

	if gs.Status.State == v1alpha1.Ready {
		logrus.WithField("gs", gs).Infof("Marking GameServer as Unhealthy")
		gsCopy := gs.DeepCopy()
		gsCopy.Status.State = v1alpha1.Unhealthy

		if _, err := hc.gameServerGetter.GameServers(gs.ObjectMeta.Namespace).Update(gsCopy); err != nil {
			return errors.Wrapf(err, "error updating GameServer %s to unhealthy", gs.ObjectMeta.Name)
		}

		hc.recorder.Event(gs, corev1.EventTypeWarning, string(gsCopy.Status.State), "GameServer container terminated")
	}

	return nil
}
