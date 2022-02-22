// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1alpha1 "github.com/HwameiStor/local-storage/pkg/apis/uds/v1alpha1"
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
	// Group=localstorage.hwameistor.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("localstoragealerts"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().LocalStorageAlerts().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("localstoragenodes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().LocalStorageNodes().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("localvolumes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().LocalVolumes().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("localvolumeexpands"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().LocalVolumeExpands().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("localvolumereplicas"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().LocalVolumeReplicas().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("physicaldisks"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Uds().V1alpha1().PhysicalDisks().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
