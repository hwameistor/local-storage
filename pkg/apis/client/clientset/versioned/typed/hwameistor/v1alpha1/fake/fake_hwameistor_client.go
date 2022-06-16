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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/hwameistor/local-storage/pkg/apis/client/clientset/versioned/typed/hwameistor/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeHwameistorV1alpha1 struct {
	*testing.Fake
}

func (c *FakeHwameistorV1alpha1) LocalStorageNodes() v1alpha1.LocalStorageNodeInterface {
	return &FakeLocalStorageNodes{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumes() v1alpha1.LocalVolumeInterface {
	return &FakeLocalVolumes{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeConverts() v1alpha1.LocalVolumeConvertInterface {
	return &FakeLocalVolumeConverts{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeExpands() v1alpha1.LocalVolumeExpandInterface {
	return &FakeLocalVolumeExpands{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeGroups() v1alpha1.LocalVolumeGroupInterface {
	return &FakeLocalVolumeGroups{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeGroupMigrates() v1alpha1.LocalVolumeGroupMigrateInterface {
	return &FakeLocalVolumeGroupMigrates{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeMigrates() v1alpha1.LocalVolumeMigrateInterface {
	return &FakeLocalVolumeMigrates{c}
}

func (c *FakeHwameistorV1alpha1) LocalVolumeReplicas() v1alpha1.LocalVolumeReplicaInterface {
	return &FakeLocalVolumeReplicas{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeHwameistorV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
