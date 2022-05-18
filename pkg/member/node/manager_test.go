package node

import (
	"fmt"
	"github.com/golang/mock/gomock"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-storage/pkg/common"
	"github.com/hwameistor/local-storage/pkg/member/node/diskmonitor"
	"github.com/hwameistor/local-storage/pkg/member/node/storage"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func Test_manager_ReconcileVolumeReplica(t *testing.T) {

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	var localVolumeReplica = &apisv1alpha1.LocalVolumeReplica{}
	localVolumeReplica.Spec.VolumeName = "volume1"
	localVolumeReplica.Spec.PoolName = "pool1"
	localVolumeReplica.Spec.NodeName = "node1"
	localVolumeReplica.Spec.RequiredCapacityBytes = 1240
	localVolumeReplica.Name = "test1"

	m := NewMockNodeManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeReplica(localVolumeReplica).
		Return(m).
		Times(1)

	m.ReconcileVolumeReplica(localVolumeReplica)

}

func Test_manager_Run(t *testing.T) {
	type fields struct {
		name                    string
		namespace               string
		apiClient               client.Client
		informersCache          cache.Cache
		replicaRecords          map[string]string
		storageMgr              *storage.LocalManager
		diskEventQueue          *diskmonitor.EventQueue
		volumeTaskQueue         *common.TaskQueue
		volumeReplicaTaskQueue  *common.TaskQueue
		localDiskClaimTaskQueue *common.TaskQueue
		localDiskTaskQueue      *common.TaskQueue
		configManager           *configManager
		logger                  *log.Entry
	}
	type args struct {
		stopCh <-chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                    tt.fields.name,
				namespace:               tt.fields.namespace,
				apiClient:               tt.fields.apiClient,
				informersCache:          tt.fields.informersCache,
				replicaRecords:          tt.fields.replicaRecords,
				storageMgr:              tt.fields.storageMgr,
				diskEventQueue:          tt.fields.diskEventQueue,
				volumeTaskQueue:         tt.fields.volumeTaskQueue,
				volumeReplicaTaskQueue:  tt.fields.volumeReplicaTaskQueue,
				localDiskClaimTaskQueue: tt.fields.localDiskClaimTaskQueue,
				localDiskTaskQueue:      tt.fields.localDiskTaskQueue,
				configManager:           tt.fields.configManager,
				logger:                  tt.fields.logger,
			}
			m.Run(tt.args.stopCh)
		})
	}
}

func Test_manager_Storage(t *testing.T) {

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockNodeManager(ctrl)
	m.
		EXPECT().
		Storage().
		Return(m).
		Times(1)

	v := m.Storage()

	fmt.Printf("Test_manager_Storage v= %+v", v)

}

func Test_manager_TakeVolumeReplicaTaskAssignment(t *testing.T) {

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	var vol = &apisv1alpha1.LocalVolume{}
	vol.Name = "vol1"
	vol.Namespace = "test1"
	vol.Spec.RequiredCapacityBytes = 1240
	vol.Spec.PoolName = "pool1"
	vol.Spec.Accessibility.Node = "node1"

	m := NewMockNodeManager(ctrl)
	m.
		EXPECT().
		TakeVolumeReplicaTaskAssignment(vol).
		Return(m).
		Times(1)

	m.TakeVolumeReplicaTaskAssignment(vol)

}
