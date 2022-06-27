package manager

import (
	"context"
	ldctr "github.com/hwameistor/local-disk-manager/pkg/controller/localdisk"
	"github.com/hwameistor/local-disk-manager/pkg/localdisk"
	"github.com/hwameistor/reliable-helper-system/pkg/apis"
	apisv1alpha1 "github.com/hwameistor/reliable-helper-system/pkg/apis/hwameistor/v1alpha1"
	migratepkg "github.com/hwameistor/reliable-helper-system/pkg/migrate"
	"github.com/hwameistor/reliable-helper-system/pkg/replacedisk/node"
	"github.com/hwameistor/reliable-helper-system/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	mgrpkg "sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	replaceDiskNodeManager apis.ReplaceDiskNodeManager
)

type manager struct {
	nodeName string

	namespace string

	apiClient client.Client

	rdhandler *ReplaceDiskHandler

	migrateCtr migratepkg.Controller

	localDiskController localdisk.Controller

	ldhandler *ldctr.LocalDiskHandler

	mgr mgrpkg.Manager

	cmdExec *lvmExecutor

	logger *log.Entry
}

// New replacedisk manager
func New(mgr mgrpkg.Manager) (apis.ReplaceDiskManager, error) {
	var recorder record.EventRecorder
	return &manager{
		nodeName:            utils.GetNodeName(),
		namespace:           utils.GetNamespace(),
		apiClient:           mgr.GetClient(),
		rdhandler:           NewReplaceDiskHandler(mgr.GetClient(), recorder),
		migrateCtr:          migratepkg.NewController(mgr),
		mgr:                 mgr,
		localDiskController: localdisk.NewController(mgr),
		ldhandler:           ldctr.NewLocalDiskHandler(mgr.GetClient(), recorder),
		cmdExec:             NewLVMExecutor(),
		logger:              log.WithField("Module", "ReplaceDisk"),
	}, nil
}

func (m *manager) Run(stopCh <-chan struct{}) {

	go m.startReplaceDiskTaskWorker(stopCh)

}

func (m *manager) ReplaceDiskNodeManager() apis.ReplaceDiskNodeManager {
	if replaceDiskNodeManager == nil {
		replaceDiskNodeManager = node.NewReplaceDiskNodeManager()
	}
	return replaceDiskNodeManager
}

// ReplaceDiskHandler
type ReplaceDiskHandler struct {
	client.Client
	record.EventRecorder
	ReplaceDisk apisv1alpha1.ReplaceDisk
}

// NewReplaceDiskHandler
func NewReplaceDiskHandler(client client.Client, recorder record.EventRecorder) *ReplaceDiskHandler {
	return &ReplaceDiskHandler{
		Client:        client,
		EventRecorder: recorder,
	}
}

// ListReplaceDisk
func (rdHandler *ReplaceDiskHandler) ListReplaceDisk() (*apisv1alpha1.ReplaceDiskList, error) {
	list := &apisv1alpha1.ReplaceDiskList{
		TypeMeta: metav1.TypeMeta{
			Kind:       ReplaceDiskKind,
			APIVersion: ReplaceDiskAPIVersion,
		},
	}

	err := rdHandler.List(context.TODO(), list)
	return list, err
}

// GetReplaceDisk
func (rdHandler *ReplaceDiskHandler) GetReplaceDisk(key client.ObjectKey) (*apisv1alpha1.ReplaceDisk, error) {
	rd := &apisv1alpha1.ReplaceDisk{}
	if err := rdHandler.Get(context.Background(), key, rd); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return rd, nil
}

// UpdateReplaceDiskStatus
func (rdHandler *ReplaceDiskHandler) UpdateReplaceDiskStatus(status apisv1alpha1.ReplaceDiskStatus) error {
	log.Debugf("ReplaceDiskHandler UpdateReplaceDiskStatus status = %v", status)

	rdHandler.ReplaceDisk.Status.OldDiskReplaceStatus = status.OldDiskReplaceStatus
	rdHandler.ReplaceDisk.Status.NewDiskReplaceStatus = status.NewDiskReplaceStatus
	return rdHandler.Status().Update(context.Background(), &rdHandler.ReplaceDisk)
}

// Refresh
func (rdHandler *ReplaceDiskHandler) Refresh() (*ReplaceDiskHandler, error) {
	rd, err := rdHandler.GetReplaceDisk(client.ObjectKey{Name: rdHandler.ReplaceDisk.GetName(), Namespace: rdHandler.ReplaceDisk.GetNamespace()})
	if err != nil {
		return rdHandler, err
	}
	rdHandler.SetReplaceDisk(*rd.DeepCopy())
	return rdHandler, nil
}

// SetReplaceDisk
func (rdHandler *ReplaceDiskHandler) SetReplaceDisk(rd apisv1alpha1.ReplaceDisk) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk = rd
	return rdHandler
}

// SetReplaceDisk
func (rdHandler *ReplaceDiskHandler) SetMigrateVolumeNames(volumeNames []string) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Status.MigrateVolumeNames = volumeNames
	return rdHandler
}

// SetMigrateSucceededVolumeNames
func (rdHandler *ReplaceDiskHandler) SetMigrateSucceededVolumeNames(volumeNames []string) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Status.MigrateSucceededVolumeNames = volumeNames
	return rdHandler
}

// SetMigrateBackUpSucceededVolumeNames
func (rdHandler *ReplaceDiskHandler) SetMigrateBackUpSucceededVolumeNames(volumeNames []string) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Status.MigrateBackUpSucceededVolumeNames = volumeNames
	return rdHandler
}

// SetMigrateFailededVolumeNames
func (rdHandler *ReplaceDiskHandler) SetMigrateFailededVolumeNames(volumeNames []string) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Status.MigrateFailededVolumeNames = volumeNames
	return rdHandler
}

// SetMigrateBackUpFailededVolumeNames
func (rdHandler *ReplaceDiskHandler) SetMigrateBackUpFailededVolumeNames(volumeNames []string) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Status.MigrateBackUpFailededVolumeNames = volumeNames
	return rdHandler
}

// SetErrMsg
func (rdHandler *ReplaceDiskHandler) SetErrMsg(errMsg string) error {
	rdHandler.ReplaceDisk.Status.ErrMsg = errMsg
	return rdHandler.Status().Update(context.Background(), &rdHandler.ReplaceDisk)
}

// SetWarnMsg
func (rdHandler *ReplaceDiskHandler) SetWarnMsg(warnMsg string) error {
	rdHandler.ReplaceDisk.Status.WarnMsg = warnMsg
	return rdHandler.Status().Update(context.Background(), &rdHandler.ReplaceDisk)
}

// ReplaceDiskStage
func (rdHandler *ReplaceDiskHandler) ReplaceDiskStage() apisv1alpha1.ReplaceDiskStage {
	return rdHandler.ReplaceDisk.Spec.ReplaceDiskStage
}

// ReplaceDiskStatus
func (rdHandler *ReplaceDiskHandler) ReplaceDiskStatus() apisv1alpha1.ReplaceDiskStatus {
	return rdHandler.ReplaceDisk.Status
}

// SetReplaceDiskStage
func (rdHandler *ReplaceDiskHandler) SetReplaceDiskStage(stage apisv1alpha1.ReplaceDiskStage) *ReplaceDiskHandler {
	rdHandler.ReplaceDisk.Spec.ReplaceDiskStage = stage
	return rdHandler
}

// UpdateReplaceDiskCR
func (rdHandler *ReplaceDiskHandler) UpdateReplaceDiskCR() error {
	return rdHandler.Update(context.Background(), &rdHandler.ReplaceDisk)
}

// SetReplaceDisk
func (rdHandler *ReplaceDiskHandler) CheckReplaceDiskTaskDestroyed(rd apisv1alpha1.ReplaceDisk) bool {
	_, err := rdHandler.GetReplaceDisk(client.ObjectKey{Name: rd.Name, Namespace: rd.Namespace})
	if err != nil {
		if errors.IsNotFound(err) {
			return true
		}
		return false
	}

	return false
}
