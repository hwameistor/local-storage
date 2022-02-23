package controller

import (
	"context"
	"strings"
	"time"

	"github.com/hwameistor/local-storage/pkg/apis"
	localstoragev1alpha1 "github.com/hwameistor/local-storage/pkg/apis/localstorage/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/sirupsen/logrus"
	coorv1 "k8s.io/api/coordination/v1"
)

const (
	nodeStatusCheckInterval = 20 * time.Second
)

func (m *manager) syncNodesStatusForever(stopCh <-chan struct{}) {
	m.logger.Debug("Starting a worker to synchronize nodes status regularly")
	m.syncNodesStatus()
	for {
		select {
		case <-time.After(nodeStatusCheckInterval):
			m.syncNodesStatus()
		case <-stopCh:
			m.logger.Debug("Exit the node status synchronizing")
			return
		}
	}
}

// syncNodesStatus is to check LocalStorageNode status regularly
// in case of a Node going offline, the status will not be set to Offline by default.
// So, this process will help to set the Node's status correctly
func (m *manager) syncNodesStatus() {
	m.logger.Debug("Checking for nodes status")
	ctx := context.TODO()

	nodeList := &localstoragev1alpha1.LocalStorageNodeList{}
	if err := m.apiClient.List(ctx, nodeList); err != nil {
		m.logger.WithError(err).Error("Failed to get NodeList")
		return
	}

	leaseList := &coorv1.LeaseList{}
	opts := []client.ListOption{
		client.InNamespace(m.namespace),
	}
	if err := m.apiClient.List(ctx, leaseList, opts...); err != nil {
		m.logger.WithError(err).Error("Failed to get LeaseList")
		return
	}
	nodeLeases := map[string]*coorv1.Lease{}
	for i, lease := range leaseList.Items {
		if !strings.HasPrefix(lease.Name, apis.NodeLeaseNamePrefix) {
			continue
		}
		if lease.Spec.HolderIdentity != nil && *lease.Spec.HolderIdentity != "" {
			nodeLeases[*lease.Spec.HolderIdentity] = &leaseList.Items[i]
		}
	}

	currTime := time.Now()
	for _, node := range nodeList.Items {
		lease, ok := nodeLeases[node.Name]
		if !ok {
			// no lease, should set node offline
			m.setNodeState(ctx, &node, localstoragev1alpha1.NodeStateOffline)
		} else if lease.Spec.LeaseDurationSeconds != nil {
			if int32(currTime.Sub(lease.Spec.RenewTime.Time).Seconds()) > *lease.Spec.LeaseDurationSeconds {
				m.setNodeState(ctx, &node, localstoragev1alpha1.NodeStateOffline)
			} else {
				m.setNodeState(ctx, &node, localstoragev1alpha1.NodeStateReady)
			}
		}
	}
}

func (m *manager) setNodeState(ctx context.Context, node *localstoragev1alpha1.LocalStorageNode, newState localstoragev1alpha1.State) error {
	logCtx := m.logger.WithFields(log.Fields{"node": node.Name, "oldState": node.Status.State, "newState": newState})
	if node.Status.State == newState {
		return nil
	}
	if node.Status.State == localstoragev1alpha1.NodeStateMaintain && newState == localstoragev1alpha1.NodeStateReady {
		return nil
	}

	node.Status.State = newState
	logCtx.Info("Updated Node state")
	return m.apiClient.Status().Update(ctx, node)
}
