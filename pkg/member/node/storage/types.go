package storage

import (
	"errors"
	"os"
	"syscall"

	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
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
//go:generate mockgen -source=types.go -destination=../../../mocks/local_pool_manager.go  -package=mocks
type LocalPoolManager interface {
	ExtendPools(localDisks []*apisv1alpha1.LocalDisk) error

	ExtendPoolsInfo(localDisks map[string]*apisv1alpha1.LocalDisk) (map[string]*apisv1alpha1.LocalPool, error)

	GetReplicas() (map[string]*apisv1alpha1.LocalVolumeReplica, error)
}

// LocalVolumeReplicaManager interface
//go:generate mockgen -source=types.go -destination=../../../mocks/local_volume_replica_manager.go  -package=mocks
type LocalVolumeReplicaManager interface {
	CreateVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) (*apisv1alpha1.LocalVolumeReplica, error)
	DeleteVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) error
	ExpandVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica, newCapacityBytes int64) (*apisv1alpha1.LocalVolumeReplica, error)
	GetVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) (*apisv1alpha1.LocalVolumeReplica, error)
	TestVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) (*apisv1alpha1.LocalVolumeReplica, error)

	// consistencyCheck on all volume replicas by comparing VolumeReplica and underlying volumes
	// will log all the check results for alerting or other purpose, but not block anything
	ConsistencyCheck()
}

// LocalDiskManager is an interface to manage local disks
//go:generate mockgen -source=types.go -destination=../../../mocks/local_disk_manager.go  -package=mocks
type LocalDiskManager interface {
	// Discover all disks including HDD, SSD, NVMe, etc..
	DiscoverAvailableDisks() ([]*apisv1alpha1.LocalDisk, error)
	GetLocalDisks() (map[string]*apisv1alpha1.LocalDisk, error)
}

// LocalRegistry interface
//go:generate mockgen -source=types.go -destination=../../../mocks/local_registry.go  -package=mocks
type LocalRegistry interface {
	Init()

	Disks() map[string]*apisv1alpha1.LocalDisk
	Pools() map[string]*apisv1alpha1.LocalPool
	VolumeReplicas() map[string]*apisv1alpha1.LocalVolumeReplica
	HasVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) bool
	UpdateNodeForVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica)
	SyncResourcesToNodeCRD(localDisks map[string]*apisv1alpha1.LocalDisk) error
}

// DeviceInfo struct
type DeviceInfo struct {
	OSFileInfo   os.FileInfo
	SysTStat     *syscall.Stat_t
	Path         string
	Name         string
	Major        uint32
	Minor        uint32
	MajMinString string
}

// LocalDeviceListInterface interface
//go:generate mockgen -source=types.go -destination=../../../mocks/local_device_list_interface.go  -package=mocks
type LocalDeviceListInterface interface {
	GetDevicesInfo(string, map[string]struct{}) map[string]*DeviceInfo
}

// LocalVolumeExecutor interface
//go:generate mockgen -source=types.go -destination=../../../mocks/local_volume_executor.go  -package=mocks
type LocalVolumeExecutor interface {
	CreateVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) (*apisv1alpha1.LocalVolumeReplica, error)
	DeleteVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) error
	ExpandVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica, newCapacityBytes int64) (*apisv1alpha1.LocalVolumeReplica, error)
	TestVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) (*apisv1alpha1.LocalVolumeReplica, error)
	// GetReplicas return all replicas
	GetReplicas() (map[string]*apisv1alpha1.LocalVolumeReplica, error)

	// consistencyCheck on all volume replicas by comparing VolumeReplica and underlying volumes
	// will log all the check results for alerting or other purpose, but not block anything
	ConsistencyCheck(crdReplicas map[string]*apisv1alpha1.LocalVolumeReplica)
}

// LocalPoolExecutor interface
//go:generate mockgen -source=types.go -destination=../../../mocks/local_pool_executor.go  -package=mocks
type LocalPoolExecutor interface {
	ExtendPools(localDisks []*apisv1alpha1.LocalDisk) error
	ExtendPoolsInfo(localDisks map[string]*apisv1alpha1.LocalDisk) (map[string]*apisv1alpha1.LocalPool, error)
	GetReplicas() (map[string]*apisv1alpha1.LocalVolumeReplica, error)
}
