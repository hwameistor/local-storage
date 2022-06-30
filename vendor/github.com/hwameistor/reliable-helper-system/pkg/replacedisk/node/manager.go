package node

import (
	"github.com/hwameistor/reliable-helper-system/pkg/apis"
	"github.com/hwameistor/reliable-helper-system/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/reliable-helper-system/pkg/common"
	"github.com/hwameistor/reliable-helper-system/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// Infinitely retry
const maxRetries = 0

type ReplaceDiskNodeMgr struct {
	nodeName string

	namespace string

	replaceDiskTaskQueue *common.TaskQueue

	logger *log.Entry
}

// New ReplaceDiskNodeManager
func NewReplaceDiskNodeManager() apis.ReplaceDiskNodeManager {
	return &ReplaceDiskNodeMgr{
		nodeName:             utils.GetNodeName(),
		namespace:            utils.GetNamespace(),
		replaceDiskTaskQueue: common.NewTaskQueue("ReplaceDisk", maxRetries),
		logger:               log.WithField("Module", "ReplaceDisk"),
	}
}

func (m *ReplaceDiskNodeMgr) ReconcileReplaceDisk(replaceDisk *v1alpha1.ReplaceDisk) {
	if replaceDisk.Spec.NodeName == m.nodeName {
		m.replaceDiskTaskQueue.Add(replaceDisk.Namespace + "/" + replaceDisk.Name)
	}
}

func (m *ReplaceDiskNodeMgr) ReplaceDiskTaskQueue() *common.TaskQueue {
	return m.replaceDiskTaskQueue
}
