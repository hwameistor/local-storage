package scheduler

import (
	"fmt"
	"github.com/golang/mock/gomock"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"testing"
)

func Test_scheduler_Allocate(t *testing.T) {

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

	m := NewMockScheduler(ctrl)
	m.
		EXPECT().
		Allocate(vol).
		Return(m).
		Times(1)

	v, err := m.Allocate(vol)

	fmt.Printf("Test_scheduler_Allocate v= %+v, err= %+v", v, err)

}

func Test_scheduler_GetNodeCandidates(t *testing.T) {

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

	m := NewMockScheduler(ctrl)
	m.
		EXPECT().
		GetNodeCandidates(vol).
		Return(m).
		Times(1)

	v, err := m.GetNodeCandidates(vol)
	fmt.Printf("Test_scheduler_GetNodeCandidates v= %+v, err= %+v", v, err)

}

func Test_scheduler_Init(t *testing.T) {

	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()

	m := NewMockScheduler(ctrl)
	m.
		EXPECT().
		Init().
		Return(m).
		Times(1)

	m.Init()
}
