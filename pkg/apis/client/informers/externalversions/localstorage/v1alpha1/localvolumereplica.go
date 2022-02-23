// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	versioned "github.com/hwameistor/local-storage/pkg/apis/client/clientset/versioned"
	internalinterfaces "github.com/hwameistor/local-storage/pkg/apis/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/hwameistor/local-storage/pkg/apis/client/listers/localstorage/v1alpha1"
	localstoragev1alpha1 "github.com/hwameistor/local-storage/pkg/apis/localstorage/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// LocalVolumeReplicaInformer provides access to a shared informer and lister for
// LocalVolumeReplicas.
type LocalVolumeReplicaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.LocalVolumeReplicaLister
}

type localVolumeReplicaInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewLocalVolumeReplicaInformer constructs a new informer for LocalVolumeReplica type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewLocalVolumeReplicaInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredLocalVolumeReplicaInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredLocalVolumeReplicaInformer constructs a new informer for LocalVolumeReplica type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredLocalVolumeReplicaInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LocalStorageV1alpha1().LocalVolumeReplicas().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LocalStorageV1alpha1().LocalVolumeReplicas().Watch(context.TODO(), options)
			},
		},
		&localstoragev1alpha1.LocalVolumeReplica{},
		resyncPeriod,
		indexers,
	)
}

func (f *localVolumeReplicaInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredLocalVolumeReplicaInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *localVolumeReplicaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&localstoragev1alpha1.LocalVolumeReplica{}, f.defaultInformer)
}

func (f *localVolumeReplicaInformer) Lister() v1alpha1.LocalVolumeReplicaLister {
	return v1alpha1.NewLocalVolumeReplicaLister(f.Informer().GetIndexer())
}
