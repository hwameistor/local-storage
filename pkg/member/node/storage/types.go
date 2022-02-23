package storage

import (
	"errors"
	localstoragev1alpha1 "github.com/hwameistor/local-storage/pkg/apis/localstorage/v1alpha1"
)

// variables
var (
	ErrorPoolNotFound                 = errors.New("not found pool")
	ErrorReplicaNotFound              = errors.New("not found replica")
	ErrorReplicaExists                = errors.New("already exists replica")
	ErrorInsufficientRequestResources = errors.New("insufficient request resources")
	ErrorOverLimitedRequestResource   = errors.New("over limited request resources")
)

// LocalPoolManager is an interface to manage local storage pools
type LocalPoolManager interface {
	ExtendPools(localDisks []*localstoragev1alpha1.LocalDisk) error

	ExtendPoolsInfo(localDisks map[string]*localstoragev1alpha1.LocalDisk) (map[string]*localstoragev1alpha1.LocalPool, error)

	GetReplicas() (map[string]*localstoragev1alpha1.LocalVolumeReplica, error)
}

// LocalVolumeReplicaManager interface
type LocalVolumeReplicaManager interface {
	CreateVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) (*localstoragev1alpha1.LocalVolumeReplica, error)
	DeleteVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) error
	ExpandVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica, newCapacityBytes int64) (*localstoragev1alpha1.LocalVolumeReplica, error)
	GetVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) (*localstoragev1alpha1.LocalVolumeReplica, error)
	TestVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) (*localstoragev1alpha1.LocalVolumeReplica, error)

	// consistencyCheck on all volume replicas by comparing VolumeReplica and underlying volumes
	// will log all the check results for alerting or other purpose, but not block anything
	ConsistencyCheck()
}

// LocalDiskManager is an interface to manage local disks
type LocalDiskManager interface {
	// Discover all disks including HDD, SSD, NVMe, etc..
	DiscoverAvailableDisks() ([]*localstoragev1alpha1.LocalDisk, error)
	GetLocalDisks() (map[string]*localstoragev1alpha1.LocalDisk, error)
}

// LocalRegistry interface
type LocalRegistry interface {
	Init()

	Disks() map[string]*localstoragev1alpha1.LocalDisk
	Pools() map[string]*localstoragev1alpha1.LocalPool
	VolumeReplicas() map[string]*localstoragev1alpha1.LocalVolumeReplica
	HasVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) bool
	UpdateNodeForVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica)
	SyncResourcesToNodeCRD(localDisks map[string]*localstoragev1alpha1.LocalDisk) error
}

// LocalDeviceListInterface interface
type LocalDeviceListInterface interface {
	GetDevicesInfo(string, map[string]struct{}) map[string]*DeviceInfo
}

// LocalVolumeExecutor interface
type LocalVolumeExecutor interface {
	CreateVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) (*localstoragev1alpha1.LocalVolumeReplica, error)
	DeleteVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) error
	ExpandVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica, newCapacityBytes int64) (*localstoragev1alpha1.LocalVolumeReplica, error)
	TestVolumeReplica(replica *localstoragev1alpha1.LocalVolumeReplica) (*localstoragev1alpha1.LocalVolumeReplica, error)
	// GetReplicas return all replicas
	GetReplicas() (map[string]*localstoragev1alpha1.LocalVolumeReplica, error)

	// consistencyCheck on all volume replicas by comparing VolumeReplica and underlying volumes
	// will log all the check results for alerting or other purpose, but not block anything
	ConsistencyCheck(crdReplicas map[string]*localstoragev1alpha1.LocalVolumeReplica)
}

// LocalPoolExecutor interface
type LocalPoolExecutor interface {
	ExtendPools(localDisks []*localstoragev1alpha1.LocalDisk) error
	ExtendPoolsInfo(localDisks map[string]*localstoragev1alpha1.LocalDisk) (map[string]*localstoragev1alpha1.LocalPool, error)
	GetReplicas() (map[string]*localstoragev1alpha1.LocalVolumeReplica, error)
}
