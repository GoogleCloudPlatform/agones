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

// This code was autogenerated. Do not edit directly.

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	allocationv1alpha1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1alpha1"
	autoscalingv1 "agones.dev/agones/pkg/client/clientset/versioned/typed/autoscaling/v1"
	multiclusterv1alpha1 "agones.dev/agones/pkg/client/clientset/versioned/typed/multicluster/v1alpha1"
	stablev1alpha1 "agones.dev/agones/pkg/client/clientset/versioned/typed/stable/v1alpha1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AllocationV1alpha1() allocationv1alpha1.AllocationV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Allocation() allocationv1alpha1.AllocationV1alpha1Interface
	AutoscalingV1() autoscalingv1.AutoscalingV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Autoscaling() autoscalingv1.AutoscalingV1Interface
	MulticlusterV1alpha1() multiclusterv1alpha1.MulticlusterV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Multicluster() multiclusterv1alpha1.MulticlusterV1alpha1Interface
	StableV1alpha1() stablev1alpha1.StableV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Stable() stablev1alpha1.StableV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	allocationV1alpha1   *allocationv1alpha1.AllocationV1alpha1Client
	autoscalingV1        *autoscalingv1.AutoscalingV1Client
	multiclusterV1alpha1 *multiclusterv1alpha1.MulticlusterV1alpha1Client
	stableV1alpha1       *stablev1alpha1.StableV1alpha1Client
}

// AllocationV1alpha1 retrieves the AllocationV1alpha1Client
func (c *Clientset) AllocationV1alpha1() allocationv1alpha1.AllocationV1alpha1Interface {
	return c.allocationV1alpha1
}

// Deprecated: Allocation retrieves the default version of AllocationClient.
// Please explicitly pick a version.
func (c *Clientset) Allocation() allocationv1alpha1.AllocationV1alpha1Interface {
	return c.allocationV1alpha1
}

// AutoscalingV1 retrieves the AutoscalingV1Client
func (c *Clientset) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	return c.autoscalingV1
}

// Deprecated: Autoscaling retrieves the default version of AutoscalingClient.
// Please explicitly pick a version.
func (c *Clientset) Autoscaling() autoscalingv1.AutoscalingV1Interface {
	return c.autoscalingV1
}

// MulticlusterV1alpha1 retrieves the MulticlusterV1alpha1Client
func (c *Clientset) MulticlusterV1alpha1() multiclusterv1alpha1.MulticlusterV1alpha1Interface {
	return c.multiclusterV1alpha1
}

// Deprecated: Multicluster retrieves the default version of MulticlusterClient.
// Please explicitly pick a version.
func (c *Clientset) Multicluster() multiclusterv1alpha1.MulticlusterV1alpha1Interface {
	return c.multiclusterV1alpha1
}

// StableV1alpha1 retrieves the StableV1alpha1Client
func (c *Clientset) StableV1alpha1() stablev1alpha1.StableV1alpha1Interface {
	return c.stableV1alpha1
}

// Deprecated: Stable retrieves the default version of StableClient.
// Please explicitly pick a version.
func (c *Clientset) Stable() stablev1alpha1.StableV1alpha1Interface {
	return c.stableV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.allocationV1alpha1, err = allocationv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.autoscalingV1, err = autoscalingv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.multiclusterV1alpha1, err = multiclusterv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.stableV1alpha1, err = stablev1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.allocationV1alpha1 = allocationv1alpha1.NewForConfigOrDie(c)
	cs.autoscalingV1 = autoscalingv1.NewForConfigOrDie(c)
	cs.multiclusterV1alpha1 = multiclusterv1alpha1.NewForConfigOrDie(c)
	cs.stableV1alpha1 = stablev1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.allocationV1alpha1 = allocationv1alpha1.New(c)
	cs.autoscalingV1 = autoscalingv1.New(c)
	cs.multiclusterV1alpha1 = multiclusterv1alpha1.New(c)
	cs.stableV1alpha1 = stablev1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
