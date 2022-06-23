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
	"context"

	v1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLocalVolumeGroupMigrates implements LocalVolumeGroupMigrateInterface
type FakeLocalVolumeGroupMigrates struct {
	Fake *FakeHwameistorV1alpha1
}

var localvolumegroupmigratesResource = schema.GroupVersionResource{Group: "hwameistor.io", Version: "v1alpha1", Resource: "localvolumegroupmigrates"}

var localvolumegroupmigratesKind = schema.GroupVersionKind{Group: "hwameistor.io", Version: "v1alpha1", Kind: "LocalVolumeGroupMigrate"}

// Get takes name of the localVolumeGroupMigrate, and returns the corresponding localVolumeGroupMigrate object, and an error if there is any.
func (c *FakeLocalVolumeGroupMigrates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LocalVolumeGroupMigrate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(localvolumegroupmigratesResource, name), &v1alpha1.LocalVolumeGroupMigrate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalVolumeGroupMigrate), err
}

// List takes label and field selectors, and returns the list of LocalVolumeGroupMigrates that match those selectors.
func (c *FakeLocalVolumeGroupMigrates) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LocalVolumeGroupMigrateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(localvolumegroupmigratesResource, localvolumegroupmigratesKind, opts), &v1alpha1.LocalVolumeGroupMigrateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.LocalVolumeGroupMigrateList{ListMeta: obj.(*v1alpha1.LocalVolumeGroupMigrateList).ListMeta}
	for _, item := range obj.(*v1alpha1.LocalVolumeGroupMigrateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested localVolumeGroupMigrates.
func (c *FakeLocalVolumeGroupMigrates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(localvolumegroupmigratesResource, opts))
}

// Create takes the representation of a localVolumeGroupMigrate and creates it.  Returns the server's representation of the localVolumeGroupMigrate, and an error, if there is any.
func (c *FakeLocalVolumeGroupMigrates) Create(ctx context.Context, localVolumeGroupMigrate *v1alpha1.LocalVolumeGroupMigrate, opts v1.CreateOptions) (result *v1alpha1.LocalVolumeGroupMigrate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(localvolumegroupmigratesResource, localVolumeGroupMigrate), &v1alpha1.LocalVolumeGroupMigrate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalVolumeGroupMigrate), err
}

// Update takes the representation of a localVolumeGroupMigrate and updates it. Returns the server's representation of the localVolumeGroupMigrate, and an error, if there is any.
func (c *FakeLocalVolumeGroupMigrates) Update(ctx context.Context, localVolumeGroupMigrate *v1alpha1.LocalVolumeGroupMigrate, opts v1.UpdateOptions) (result *v1alpha1.LocalVolumeGroupMigrate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(localvolumegroupmigratesResource, localVolumeGroupMigrate), &v1alpha1.LocalVolumeGroupMigrate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalVolumeGroupMigrate), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLocalVolumeGroupMigrates) UpdateStatus(ctx context.Context, localVolumeGroupMigrate *v1alpha1.LocalVolumeGroupMigrate, opts v1.UpdateOptions) (*v1alpha1.LocalVolumeGroupMigrate, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(localvolumegroupmigratesResource, "status", localVolumeGroupMigrate), &v1alpha1.LocalVolumeGroupMigrate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalVolumeGroupMigrate), err
}

// Delete takes name of the localVolumeGroupMigrate and deletes it. Returns an error if one occurs.
func (c *FakeLocalVolumeGroupMigrates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(localvolumegroupmigratesResource, name), &v1alpha1.LocalVolumeGroupMigrate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLocalVolumeGroupMigrates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(localvolumegroupmigratesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.LocalVolumeGroupMigrateList{})
	return err
}

// Patch applies the patch and returns the patched localVolumeGroupMigrate.
func (c *FakeLocalVolumeGroupMigrates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LocalVolumeGroupMigrate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(localvolumegroupmigratesResource, name, pt, data, subresources...), &v1alpha1.LocalVolumeGroupMigrate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalVolumeGroupMigrate), err
}
