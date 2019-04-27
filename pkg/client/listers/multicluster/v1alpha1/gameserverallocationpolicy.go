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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "agones.dev/agones/pkg/apis/multicluster/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GameServerAllocationPolicyLister helps list GameServerAllocationPolicies.
type GameServerAllocationPolicyLister interface {
	// List lists all GameServerAllocationPolicies in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GameServerAllocationPolicy, err error)
	// GameServerAllocationPolicies returns an object that can list and get GameServerAllocationPolicies.
	GameServerAllocationPolicies(namespace string) GameServerAllocationPolicyNamespaceLister
	GameServerAllocationPolicyListerExpansion
}

// gameServerAllocationPolicyLister implements the GameServerAllocationPolicyLister interface.
type gameServerAllocationPolicyLister struct {
	indexer cache.Indexer
}

// NewGameServerAllocationPolicyLister returns a new GameServerAllocationPolicyLister.
func NewGameServerAllocationPolicyLister(indexer cache.Indexer) GameServerAllocationPolicyLister {
	return &gameServerAllocationPolicyLister{indexer: indexer}
}

// List lists all GameServerAllocationPolicies in the indexer.
func (s *gameServerAllocationPolicyLister) List(selector labels.Selector) (ret []*v1alpha1.GameServerAllocationPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GameServerAllocationPolicy))
	})
	return ret, err
}

// GameServerAllocationPolicies returns an object that can list and get GameServerAllocationPolicies.
func (s *gameServerAllocationPolicyLister) GameServerAllocationPolicies(namespace string) GameServerAllocationPolicyNamespaceLister {
	return gameServerAllocationPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GameServerAllocationPolicyNamespaceLister helps list and get GameServerAllocationPolicies.
type GameServerAllocationPolicyNamespaceLister interface {
	// List lists all GameServerAllocationPolicies in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GameServerAllocationPolicy, err error)
	// Get retrieves the GameServerAllocationPolicy from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GameServerAllocationPolicy, error)
	GameServerAllocationPolicyNamespaceListerExpansion
}

// gameServerAllocationPolicyNamespaceLister implements the GameServerAllocationPolicyNamespaceLister
// interface.
type gameServerAllocationPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GameServerAllocationPolicies in the indexer for a given namespace.
func (s gameServerAllocationPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GameServerAllocationPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GameServerAllocationPolicy))
	})
	return ret, err
}

// Get retrieves the GameServerAllocationPolicy from the indexer for a given namespace and name.
func (s gameServerAllocationPolicyNamespaceLister) Get(name string) (*v1alpha1.GameServerAllocationPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("gameserverallocationpolicy"), name)
	}
	return obj.(*v1alpha1.GameServerAllocationPolicy), nil
}
