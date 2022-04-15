package localdisk

import (
	"context"
	ldm "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-disk-manager/pkg/disk/manager"
	"github.com/hwameistor/local-disk-manager/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

// Controller The smallest unit to be processed here should be the disk.
// The main thing to do is how to convert the local disk into resources in the cluster
type Controller struct {
	// mgr k8s runtime controller
	mgr crmanager.Manager

	// Namespace is the namespace in which ldm is installed
	NameSpace string

	// NodeName is the node in which ldm is installed
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

// CreateLocalDisk
func (ctr Controller) CreateLocalDisk(ld ldm.LocalDisk) error {
	log.Debugf("Create LocalDisk for %+v", ld)
	return ctr.mgr.GetClient().Create(context.Background(), &ld)
}

// CreateLocalDisk
func (ctr Controller) UpdateLocalDisk(ld ldm.LocalDisk) error {
	newLd := ld.DeepCopy()
	key := client.ObjectKey{Name: ld.GetName(), Namespace: ""}

	oldLd, err := ctr.GetLocalDisk(key)
	if err != nil {
		return err
	}

	// TODO: merge old disk and new disk
	ctr.mergerLocalDisk(oldLd, newLd)
	return ctr.mgr.GetClient().Update(context.Background(), newLd)
}

// IsAlreadyExist
func (ctr Controller) IsAlreadyExist(ld ldm.LocalDisk) bool {
	key := client.ObjectKey{Name: ld.GetName(), Namespace: ""}
	if lookLd, err := ctr.GetLocalDisk(key); err != nil {
		return false
	} else {
		return ld.GetName() == lookLd.GetName()
	}
}

// GetLocalDisk
func (ctr Controller) GetLocalDisk(key client.ObjectKey) (ldm.LocalDisk, error) {
	ld := ldm.LocalDisk{}
	if err := ctr.mgr.GetClient().Get(context.Background(), key, &ld); err != nil {
		if errors.IsNotFound(err) {
			return ld, nil
		}
		return ld, err
	}

	return ld, nil
}

// ConvertDiskToLocalDisk
func (ctr Controller) ConvertDiskToLocalDisk(disk manager.DiskInfo) (ld ldm.LocalDisk) {
	ld.TypeMeta = metav1.TypeMeta{
		Kind:       LocalDiskKind,
		APIVersion: LocalDiskAPIVersion,
	}
	ld.ObjectMeta = metav1.ObjectMeta{
		Namespace:   ctr.NameSpace,
		Labels:      make(map[string]string),
		Annotations: make(map[string]string),
		Name:        ctr.GenLocalDiskName(disk),
	}

	ld.Spec = ctr.genLocalDiskSpec(disk)
	ld.Status = ctr.genLocalDiskStatus()
	return
}

// GetPartitionInfo
func (ctr Controller) GetPartitionInfo(originParts []manager.PartitionInfo) (ldmParts []ldm.PartitionInfo) {
	for _, part := range originParts {
		ldmPart := ldm.PartitionInfo{}
		ldmPart.HasFileSystem = true
		ldmPart.FileSystem.Type = part.Filesystem
		ldmParts = append(ldmParts, ldmPart)
	}
	return
}

// genLocalDiskSpec
func (ctr Controller) genLocalDiskSpec(disk manager.DiskInfo) ldm.LocalDiskSpec {
	spec := ldm.LocalDiskSpec{State: ldm.LocalDiskActive}
	spec.NodeName = utils.ConvertNodeName(ctr.NodeName)

	spec.Capacity = disk.Attribute.Capacity
	spec.DevicePath = disk.Attribute.DevName

	spec.DiskAttributes.Type = disk.Attribute.DriverType
	spec.DiskAttributes.Vendor = disk.Attribute.Vendor
	spec.DiskAttributes.ModelName = disk.Attribute.Model
	spec.DiskAttributes.Protocol = disk.Attribute.Bus
	spec.DiskAttributes.SerialNumber = disk.Attribute.Serial
	spec.DiskAttributes.DevType = disk.Attribute.DevType

	spec.UUID = disk.GenerateUUID()
	spec.PartitionInfo = ctr.GetPartitionInfo(disk.Partitions)

	return spec
}

func (ctr Controller) mergerLocalDisk(oldLd ldm.LocalDisk, newLd *ldm.LocalDisk) {
	newLd.Status = oldLd.Status
	newLd.TypeMeta = oldLd.TypeMeta
	newLd.ObjectMeta = oldLd.ObjectMeta
}

// genLocalDiskStatus
func (ctr Controller) genLocalDiskStatus() ldm.LocalDiskStatus {
	return ldm.LocalDiskStatus{State: ldm.LocalDiskUnclaimed}
}

func (ctr Controller) GenLocalDiskName(disk manager.DiskInfo) string {
	return utils.ConvertNodeName(ctr.NodeName) + "-" + disk.Name
}
