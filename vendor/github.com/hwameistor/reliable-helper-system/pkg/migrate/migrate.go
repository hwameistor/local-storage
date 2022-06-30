package migrate

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	lsapisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/reliable-helper-system/pkg/utils"
	log "github.com/sirupsen/logrus"
	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

// Controller
type Controller struct {
	// mgr k8s runtime controller
	mgr crmanager.Manager

	// Namespace is the namespace in which migration is created
	NameSpace string

	// NodeName is the node in which migration is created
	NodeName string
}

// NewController
func NewController(mgr crmanager.Manager) Controller {
	return Controller{
		mgr:       mgr,
		NameSpace: utils.GetNamespace(),
		NodeName:  utils.GetNodeName(),
	}
}

// CreateLocalVolumeMigrate
func (ctr Controller) CreateLocalVolumeMigrate(lvm lsapisv1alpha1.LocalVolumeMigrate) error {
	log.Debugf("Create LocalVolumeMigrate lvm %+v", lvm)

	return ctr.mgr.GetClient().Create(context.Background(), &lvm)
}

// UpdateLocalVolumeMigrate
func (ctr Controller) UpdateLocalVolumeMigrate(lvm lsapisv1alpha1.LocalVolumeMigrate) error {
	return ctr.mgr.GetClient().Update(context.Background(), lvm.DeepCopy())
}

// IsAlreadyExist
func (ctr Controller) IsAlreadyExist(lvm lsapisv1alpha1.LocalVolumeMigrate) bool {
	key := client.ObjectKey{Name: lvm.GetName(), Namespace: lvm.Namespace}
	if lookLd, err := ctr.GetLocalVolumeMigrate(key); err != nil {
		return false
	} else {
		return lvm.GetName() == lookLd.GetName()
	}
}

// GetLocalVolumeMigrate
func (ctr Controller) GetLocalVolumeMigrate(key client.ObjectKey) (lsapisv1alpha1.LocalVolumeMigrate, error) {
	log.Debugf("GetLocalVolumeMigrate key = %+v", key)

	lvm := lsapisv1alpha1.LocalVolumeMigrate{}
	if err := ctr.mgr.GetClient().Get(context.Background(), key, &lvm); err != nil {
		if errors.IsNotFound(err) {
			return lvm, nil
		}
		return lvm, err
	}

	return lvm, nil
}

func (ctr Controller) GetLocalVolumeMigrateStatusByLocalVolume(vol lsapisv1alpha1.LocalVolume) (lsapisv1alpha1.LocalVolumeMigrateStatus, error) {
	log.Debugf("GetLocalVolumeMigrateStatusByLocalVolume vol = %+v", vol)

	var localVolumeMigrateStatus lsapisv1alpha1.LocalVolumeMigrateStatus

	migrateName := ctr.genLocalVolumeMigrateName(vol)
	key := client.ObjectKey{Name: migrateName, Namespace: vol.Namespace}
	lvm, err := ctr.GetLocalVolumeMigrate(key)
	if err != nil {
		log.Errorf("Failed to GetLocalVolumeMigrate")
		return localVolumeMigrateStatus, err
	}
	localVolumeMigrateStatus = lvm.Status

	return localVolumeMigrateStatus, nil
}

// GetLocalVolumeMigrateStatus
func (ctr Controller) GetLocalVolumeMigrateStatus(key client.ObjectKey) (lsapisv1alpha1.LocalVolumeMigrateStatus, error) {
	log.Debugf("GetLocalVolumeMigrateStatus key = %+v", key)

	var localVolumeMigrateStatus lsapisv1alpha1.LocalVolumeMigrateStatus
	lvm, err := ctr.GetLocalVolumeMigrate(key)
	if err != nil {
		log.Errorf("Failed to GetLocalVolumeMigrateStatus")
		return localVolumeMigrateStatus, err
	}
	localVolumeMigrateStatus = lvm.Status

	return localVolumeMigrateStatus, nil
}

// ConstructLocalVolumeMigrate
func (ctr Controller) ConstructLocalVolumeMigrate(vol lsapisv1alpha1.LocalVolume) (*lsapisv1alpha1.LocalVolumeMigrate, error) {
	log.Debugf("ConstructLocalVolumeMigrate vol = %+v", vol)
	selectedMigratedNodeName, err := ctr.genLocalVolumeNodeName(vol)
	if err != nil {
		log.Errorf("Failed to genLocalVolumeNodeName")
		return nil, err
	}

	var lvm = &lsapisv1alpha1.LocalVolumeMigrate{}

	lvm.TypeMeta = metav1.TypeMeta{
		Kind:       LocalVolumeMigrateKind,
		APIVersion: LocalVolumeMigrateAPIVersion,
	}
	lvm.ObjectMeta = metav1.ObjectMeta{
		Namespace:   ctr.NameSpace,
		Labels:      make(map[string]string),
		Annotations: make(map[string]string),
		Name:        ctr.genLocalVolumeMigrateName(vol),
	}

	lvm.Spec = ctr.genLocalVolumeMigrateSpec(lvm, selectedMigratedNodeName)

	return lvm, nil
}

// GenLocalVolumeMigrateNodeName
func (ctr Controller) genLocalVolumeNodeName(vol lsapisv1alpha1.LocalVolume) (string, error) {
	log.Debugf("genLocalVolumeNodeName vol = %+v", vol)
	nodeList := &lsapisv1alpha1.LocalStorageNodeList{}
	if err := ctr.mgr.GetClient().List(context.Background(), nodeList); err != nil {
		log.Errorf("Failed to get NodeList")
		return "", err
	}

	var selectedMigrateNode string
	var maxFreeCapacityBytesNodes int64
	for _, node := range nodeList.Items {
		for _, affinityNode := range vol.Spec.Accessibility.Nodes {
			if node.Name == affinityNode {
				continue
			}
		}
		tmpFreeCapacityBytesNodes := node.Status.Pools[vol.Spec.PoolName].FreeCapacityBytes
		if tmpFreeCapacityBytesNodes > maxFreeCapacityBytesNodes {
			maxFreeCapacityBytesNodes = tmpFreeCapacityBytesNodes
			selectedMigrateNode = node.Name
		}
	}

	return selectedMigrateNode, nil
}

func (ctr Controller) genLocalVolumeMigrateSpec(lvm *lsapisv1alpha1.LocalVolumeMigrate, nodeName string) lsapisv1alpha1.LocalVolumeMigrateSpec {
	log.Debugf("genLocalVolumeMigrateSpec lvm = %+v, nodeName = %+v", lvm, nodeName)

	var localVolumeMigrate = lsapisv1alpha1.LocalVolumeMigrateSpec{}
	splits := strings.Split(lvm.Name, "-")
	var localVolumeName = "default"
	if len(splits) >= 2 {
		localVolumeName = strings.Split(lvm.Name, "--")[1]
	}
	localVolumeMigrate.VolumeName = localVolumeName
	localVolumeMigrate.NodeName = nodeName
	return localVolumeMigrate
}

func (ctr Controller) genLocalVolumeMigrateName(vol lsapisv1alpha1.LocalVolume) string {
	if len(vol.Spec.Accessibility.Nodes) == 0 {
		return ctr.NodeName + "--" + vol.Name
	}
	if len(vol.Spec.Accessibility.Nodes) == 1 {
		return vol.Spec.Accessibility.Nodes[0] + "--" + vol.Name
	}
	return vol.Spec.Accessibility.Nodes[0] + "--" + vol.Name
}
