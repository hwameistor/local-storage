/*
Copyright 2022 Contributors to The HwameiStor project.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// LocalVolumeReplicaLister helps list LocalVolumeReplicas.
type LocalVolumeReplicaLister interface {
	// List lists all LocalVolumeReplicas in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.LocalVolumeReplica, err error)
	// Get retrieves the LocalVolumeReplica from the index for a given name.
	Get(name string) (*v1alpha1.LocalVolumeReplica, error)
	LocalVolumeReplicaListerExpansion
}

// localVolumeReplicaLister implements the LocalVolumeReplicaLister interface.
type localVolumeReplicaLister struct {
	indexer cache.Indexer
}

// NewLocalVolumeReplicaLister returns a new LocalVolumeReplicaLister.
func NewLocalVolumeReplicaLister(indexer cache.Indexer) LocalVolumeReplicaLister {
	return &localVolumeReplicaLister{indexer: indexer}
}

// List lists all LocalVolumeReplicas in the indexer.
func (s *localVolumeReplicaLister) List(selector labels.Selector) (ret []*v1alpha1.LocalVolumeReplica, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.LocalVolumeReplica))
	})
	return ret, err
}

// Get retrieves the LocalVolumeReplica from the index for a given name.
func (s *localVolumeReplicaLister) Get(name string) (*v1alpha1.LocalVolumeReplica, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("localvolumereplica"), name)
	}
	return obj.(*v1alpha1.LocalVolumeReplica), nil
}
