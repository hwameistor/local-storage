package apis

import (
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	nodestorage "github.com/hwameistor/local-storage/pkg/member/node/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// consts
const (
	Version = "1.0.0"

	NodeLeaseNamePrefix = "hwameistor-local-storage-worker"

	ControllerLeaseName = "hwameistor-local-storage-master"

	TopologyNodeKey = "topology.lvm.hwameistor.io/node"
)

// LocalStorageMember interface
//go:generate mockgen -source=member.go -destination=../mocks/local_storage_member.go  -package=mocks
type LocalStorageMember interface {
	Run(stopCh <-chan struct{})

	// ******  configuration ******* //
	ConfigureBase(name string, namespace string, haSystemConfig apisv1alpha1.SystemConfig, cli client.Client, informersCache cache.Cache) LocalStorageMember

	ConfigureNode() LocalStorageMember

	ConfigureController(scheme *runtime.Scheme) LocalStorageMember

	ConfigureCSIDriver(driverName string, sockAddr string) LocalStorageMember

	ConfigureRESTServer(httpPort int) LocalStorageMember

	// access the modules
	Controller() ControllerManager

	Node() NodeManager

	Name() string

	Version() string

	DriverName() string
}

// ControllerManager interface
//go:generate mockgen -source=member.go -destination=../mocks/controller_manager.go  -package=mocks
type ControllerManager interface {
	Run(stopCh <-chan struct{})

	ReconcileNode(node *apisv1alpha1.LocalStorageNode)

	ReconcileVolume(vol *apisv1alpha1.LocalVolume)

	ReconcileVolumeGroup(vol *apisv1alpha1.LocalVolumeGroup)

	ReconcileVolumeExpand(vol *apisv1alpha1.LocalVolumeExpand)

	ReconcileVolumeMigrate(vol *apisv1alpha1.LocalVolumeMigrate)

	ReconcileVolumeConvert(vol *apisv1alpha1.LocalVolumeConvert)
}

// NodeManager interface
//go:generate mockgen -source=member.go -destination=../mocks/node_manager.go  -package=mocks
type NodeManager interface {
	Run(stopCh <-chan struct{})

	Storage() *nodestorage.LocalManager

	TakeVolumeReplicaTaskAssignment(vol *apisv1alpha1.LocalVolume)

	ReconcileVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica)
}
