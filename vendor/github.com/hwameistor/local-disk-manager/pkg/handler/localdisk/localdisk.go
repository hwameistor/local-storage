package localdisk

import (
	"context"

	ldm "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-disk-manager/pkg/filter"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LocalDiskHandler
type LocalDiskHandler struct {
	client.Client
	record.EventRecorder
	Ld     ldm.LocalDisk
	filter filter.LocalDiskFilter
}

// NewLocalDiskHandler
func NewLocalDiskHandler(client client.Client, recorder record.EventRecorder) *LocalDiskHandler {
	return &LocalDiskHandler{
		Client:        client,
		EventRecorder: recorder,
	}
}

// GetLocalDisk
func (ldHandler *LocalDiskHandler) GetLocalDisk(key client.ObjectKey) (*ldm.LocalDisk, error) {
	ld := ldm.LocalDisk{}
	if err := ldHandler.Get(context.Background(), key, &ld); err != nil {
		return nil, err
	}

	return &ld, nil
}

func (ldHandler *LocalDiskHandler) GetLocalDiskWithLabels(labels labels.Set) (*ldm.LocalDiskList, error) {
	list := &ldm.LocalDiskList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LocalDisk",
			APIVersion: "v1alpha1",
		},
	}
	return list, ldHandler.List(context.TODO(), list, &client.ListOptions{LabelSelector: labels.AsSelector()})
}

// ListLocalDisk
func (ldHandler *LocalDiskHandler) ListLocalDisk() (*ldm.LocalDiskList, error) {
	list := &ldm.LocalDiskList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LocalDisk",
			APIVersion: "v1alpha1",
		},
	}

	err := ldHandler.List(context.TODO(), list)
	return list, err
}

// ListNodeLocalDisk
func (ldHandler *LocalDiskHandler) ListNodeLocalDisk(node string) (*ldm.LocalDiskList, error) {
	list := &ldm.LocalDiskList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LocalDisk",
			APIVersion: "v1alpha1",
		},
	}
	nodeMatcher := client.MatchingFields{"spec.nodeName": node}
	err := ldHandler.List(context.TODO(), list, nodeMatcher)
	return list, err
}

// For
func (ldHandler *LocalDiskHandler) For(ld ldm.LocalDisk) *LocalDiskHandler {
	ldHandler.Ld = ld
	ldHandler.filter = filter.NewLocalDiskFilter(ld)
	return ldHandler
}

// UnClaimed Bounded
func (ldHandler *LocalDiskHandler) UnClaimed() bool {
	return !ldHandler.filter.
		Init().
		Unclaimed().
		GetTotalResult()
}

// BoundTo assign disk to ldc
func (ldHandler *LocalDiskHandler) BoundTo(ldc ldm.LocalDiskClaim) error {
	ldcRef, err := reference.GetReference(nil, &ldc)
	if err != nil {
		return err
	}

	ldHandler.Ld.Spec.ClaimRef = ldcRef
	ldHandler.Ld.Status.State = ldm.LocalDiskClaimed

	if err = ldHandler.UpdateStatus(); err != nil {
		return err
	}
	ldHandler.EventRecorder.Eventf(&ldHandler.Ld, v1.EventTypeNormal, "LocalDiskClaimed", "Claimed by %v/%v", ldc.Namespace, ldc.Name)
	return nil
}

// UpdateStatus
func (ldHandler *LocalDiskHandler) SetupStatus(status ldm.LocalDiskClaimState) {
	ldHandler.Ld.Status.State = status
}

// SetupLabel
func (ldHandler *LocalDiskHandler) SetupLabel(labels labels.Set) {
	if ldHandler.Ld.ObjectMeta.Labels == nil {
		ldHandler.Ld.ObjectMeta.Labels = make(map[string]string)
	}
	for k, v := range labels {
		ldHandler.Ld.ObjectMeta.Labels[k] = v
	}
}

// SetupLabel
func (ldHandler *LocalDiskHandler) RemoveLabel(labels labels.Set) {
	for k := range labels {
		delete(ldHandler.Ld.ObjectMeta.Labels, k)
	}
}

// UpdateStatus
func (ldHandler *LocalDiskHandler) UpdateStatus() error {
	return ldHandler.Update(context.Background(), &ldHandler.Ld)
}

// ClaimRef
func (ldHandler *LocalDiskHandler) ClaimRef() *v1.ObjectReference {
	return ldHandler.Ld.Spec.ClaimRef
}

// FilterDisk
func (ldHandler *LocalDiskHandler) FilterDisk(ldc ldm.LocalDiskClaim) bool {
	return ldHandler.filter.
		Init().
		Unclaimed().
		NodeMatch(ldc.Spec.NodeName).
		Capacity(ldc.Spec.Description.Capacity).
		DiskType(ldc.Spec.Description.DiskType).
		Unique(ldc.Spec.DiskRefs).
		DevType().
		NoPartition().
		GetTotalResult()
}
