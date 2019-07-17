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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	autoscaling_v1 "agones.dev/agones/pkg/apis/autoscaling/v1"
	v1alpha1 "agones.dev/agones/pkg/apis/multicluster/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=agones.dev, Version=v1
	case v1.SchemeGroupVersion.WithResource("fleets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Agones().V1().Fleets().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("gameservers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Agones().V1().GameServers().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("gameserversets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Agones().V1().GameServerSets().Informer()}, nil

		// Group=autoscaling.agones.dev, Version=v1
	case autoscaling_v1.SchemeGroupVersion.WithResource("fleetautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Autoscaling().V1().FleetAutoscalers().Informer()}, nil

		// Group=multicluster, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("gameserverallocationpolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Multicluster().V1alpha1().GameServerAllocationPolicies().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
