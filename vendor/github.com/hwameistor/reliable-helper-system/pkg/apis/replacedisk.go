package apis

import (
	apisv1alpha1 "github.com/hwameistor/reliable-helper-system/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/reliable-helper-system/pkg/common"
)

// ReplaceDiskManager interface
type ReplaceDiskManager interface {
	Run(stopCh <-chan struct{})
	ReplaceDiskNodeManager() ReplaceDiskNodeManager
}

// ReplaceDiskNodeManager interface
type ReplaceDiskNodeManager interface {
	ReconcileReplaceDisk(replaceDisk *apisv1alpha1.ReplaceDisk)
	ReplaceDiskTaskQueue() *common.TaskQueue
}
