package controller

import (
	"github.com/golang/mock/gomock"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"testing"
)

func Test_manager_ReconcileNode(t *testing.T) {
	var node = &apisv1alpha1.LocalStorageNode{}
	node.Name = "test_node1"
	node.Namespace = "test"
	node.Spec.HostName = "test_host_name1"
	node.Spec.StorageIP = "127.0.0.1"

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileNode(node).
		Return().
		Times(1)

	m.ReconcileNode(node)
}

func Test_manager_ReconcileVolume(t *testing.T) {
	var vol = &apisv1alpha1.LocalVolume{}
	vol.Name = "test_lv1"
	vol.Namespace = "test"
	vol.Spec.VolumeGroup = "test_vg1"
	vol.Spec.PersistentVolumeClaimName = "test_pvc1"
	vol.Spec.ReplicaNumber = 1
	vol.Spec.Convertible = true
	vol.Spec.PersistentVolumeClaimNamespace = "test"

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolume(vol).
		Return().
		Times(1)

	m.ReconcileVolume(vol)
}

func Test_manager_ReconcileVolumeConvert(t *testing.T) {
	var convert = &apisv1alpha1.LocalVolumeConvert{}
	convert.Name = "test_lvc1"
	convert.Namespace = "test"
	convert.Spec.VolumeName = "test_vol1"
	convert.Spec.ReplicaNumber = 1

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeConvert(convert).
		Return().
		Times(1)

	m.ReconcileVolumeConvert(convert)
}

func Test_manager_ReconcileVolumeExpand(t *testing.T) {
	var expand = &apisv1alpha1.LocalVolumeExpand{}
	expand.Name = "test_lve1"
	expand.Namespace = "test"
	expand.Spec.VolumeName = "test_vol1"
	expand.Spec.RequiredCapacityBytes = 102400

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeExpand(expand).
		Return().
		Times(1)

	m.ReconcileVolumeExpand(expand)
}

func Test_manager_ReconcileVolumeGroup(t *testing.T) {
	var volGroup = &apisv1alpha1.LocalVolumeGroup{}
	volGroup.Name = "test_vg1"
	volGroup.Namespace = "test"
	volGroup.Spec.Pods = []string{"test_pod1"}
	volGroup.Spec.Accessibility.Nodes = []string{"node1"}
	volGroup.Spec.Volumes = []apisv1alpha1.VolumeInfo{{LocalVolumeName: "local-volume-test1", PersistentVolumeClaimName: "pvc-test1"}}

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeGroup(volGroup).
		Return().
		Times(1)

	m.ReconcileVolumeGroup(volGroup)
}

func Test_manager_ReconcileVolumeGroupConvert(t *testing.T) {
	var lvgconvert = &apisv1alpha1.LocalVolumeGroupConvert{}
	lvgconvert.Name = "test_lvgc1"
	lvgconvert.Namespace = "test"
	lvgconvert.Spec.LocalVolumeGroupName = "test_lvg1"
	lvgconvert.Spec.ReplicaNumber = 2

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeGroupConvert(lvgconvert).
		Return().
		Times(1)

	m.ReconcileVolumeGroupConvert(lvgconvert)
}

func Test_manager_ReconcileVolumeGroupMigrate(t *testing.T) {
	var lvgmigrate = &apisv1alpha1.LocalVolumeGroupMigrate{}
	lvgmigrate.Name = "test_lvgm1"
	lvgmigrate.Namespace = "test"
	lvgmigrate.Spec.TargetNodesNames = []string{"test_node1"}
	lvgmigrate.Spec.SourceNodesNames = []string{"test_node2"}
	lvgmigrate.Spec.LocalVolumeGroupName = "test_lvg1"

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeGroupMigrate(lvgmigrate).
		Return().
		Times(1)

	m.ReconcileVolumeGroupMigrate(lvgmigrate)
}

func Test_manager_ReconcileVolumeMigrate(t *testing.T) {
	var migrate = &apisv1alpha1.LocalVolumeMigrate{}
	migrate.Name = "test_lvm1"
	migrate.Namespace = "test"
	migrate.Spec.VolumeName = "test_vol1"
	migrate.Spec.NodeName = "test_node"

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		ReconcileVolumeMigrate(migrate).
		Return().
		Times(1)

	m.ReconcileVolumeMigrate(migrate)
}

func Test_manager_Run(t *testing.T) {
	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	var stopCh <-chan struct{}

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		Run(stopCh).
		Return().
		Times(1)

	m.Run(stopCh)
}

func Test_manager_VolumeGroupManager(t *testing.T) {
	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		VolumeGroupManager().
		Return(m).
		Times(1)

	m.VolumeGroupManager()
}

func Test_manager_VolumeScheduler(t *testing.T) {
	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()
	var volScheduler apisv1alpha1.VolumeScheduler

	m := NewMockControllerManager(ctrl)
	m.
		EXPECT().
		VolumeScheduler().
		Return(volScheduler).
		Times(1)

	m.VolumeScheduler()
}
