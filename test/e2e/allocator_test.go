// Copyright 2019 Google LLC All Rights Reserved.
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

package e2e

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"agones.dev/agones/pkg/apis/allocation/v1alpha1"
	stablev1alpha1 "agones.dev/agones/pkg/apis/stable/v1alpha1"
	e2e "agones.dev/agones/test/e2e/framework"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAllocator(t *testing.T) {
	t.Parallel()

	kubeCore := framework.KubeClient.CoreV1()
	svc, err := kubeCore.Services("agones-system").Get("gameserver-allocator", metav1.GetOptions{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, svc.Status.LoadBalancer) {
		return
	}
	if !assert.Equal(t, 1, len(svc.Status.LoadBalancer.Ingress)) {
		return
	}
	if !assert.NotNil(t, 0, svc.Status.LoadBalancer.Ingress[0].IP) {
		return
	}

	port := svc.Spec.Ports[0]
	requestURL := fmt.Sprintf("https://%s:%d/v1alpha1/gameserverallocation", svc.Status.LoadBalancer.Ingress[0].IP, port.Port)

	flt, err := createFleet()
	if !assert.Nil(t, err) {
		return
	}
	framework.WaitForFleetCondition(t, flt, e2e.FleetReadyCount(flt.Spec.Replicas))
	gsa := &v1alpha1.GameServerAllocation{
		Spec: v1alpha1.GameServerAllocationSpec{
			Required: metav1.LabelSelector{MatchLabels: map[string]string{stablev1alpha1.FleetNameLabel: flt.ObjectMeta.Name}},
		}}

	body, err := json.Marshal(gsa)
	if !assert.Nil(t, err) {
		return
	}

	client, err := creatRestClient("agones-system", "allocator-tls")
	if !assert.Nil(t, err) {
		return
	}
	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if !assert.Nil(t, err) {
		return
	}
	defer response.Body.Close() // nolint: errcheck

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err = ioutil.ReadAll(response.Body)
	if !assert.Nil(t, err) {
		return
	}
	result := v1alpha1.GameServerAllocation{}
	err = json.Unmarshal(body, &result)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, v1alpha1.GameServerAllocationAllocated, result.Status.State)
}

// creatRestClient creates a rest client with proper certs to make a remote call.
func creatRestClient(namespace, clientSecretName string) (*http.Client, error) {
	kubeCore := framework.KubeClient.CoreV1()
	clientSecret, err := kubeCore.Secrets(namespace).Get(clientSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Create http client using cert
	clientCert := clientSecret.Data["tls.crt"]
	clientKey := clientSecret.Data["tls.key"]
	if clientCert == nil || clientKey == nil {
		return nil, fmt.Errorf("missing certificate")
	}

	// Load client cert
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	// Setup HTTPS client
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				GetClientCertificate: func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return &cert, nil
				},
			},
		},
	}, nil
}

func createFleet() (*stablev1alpha1.Fleet, error) {
	fleets := framework.AgonesClient.StableV1alpha1().Fleets(defaultNs)
	fleet := defaultFleet()
	return fleets.Create(fleet)
}
