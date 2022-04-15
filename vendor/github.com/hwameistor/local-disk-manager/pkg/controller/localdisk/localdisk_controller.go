package localdisk

import (
	"context"
	"github.com/hwameistor/local-disk-manager/pkg/filter"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"

	hwameistorv1alpha1 "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	ldmv1alpha1 "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new LocalDisk Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLocalDisk{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("localdisk-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("localdisk-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource LocalDisk
	err = c.Watch(&source.Kind{Type: &hwameistorv1alpha1.LocalDisk{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLocalDisk implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLocalDisk{}

// ReconcileLocalDisk reconciles a LocalDisk object
type ReconcileLocalDisk struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a LocalDisk object and makes changes based on the state read
// and what is in the LocalDisk.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLocalDisk) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.Infof("Reconcile LocalDisk %s", req.Name)

	ldHandler := NewLocalDiskHandler(r.Client, r.Recorder)
	ld, err := ldHandler.GetLocalDisk(req.NamespacedName)
	if err != nil {
		log.WithError(err).Errorf("Failed to get localdisk")
		return reconcile.Result{}, err
	}

	if ld != nil {
		ldHandler.For(*ld.DeepCopy())
	} else {
		// Not found
		return reconcile.Result{}, nil
	}

	// NOTE: The control logic of localdisk should only respond to the update events of the disk itself.
	// For example, if the disk smart check fails, it should update its health status at this time.
	// As for the upper layer resources that depend on it, such as LDC, what it should do is to monitor
	// the event changes of LD and adjust the changed contents accordingly.
	// The connection between them should only be related to the state.

	// Update status
	if ldHandler.ClaimRef() != nil && ldHandler.UnClaimed() {
		ldHandler.SetupStatus(ldmv1alpha1.LocalDiskClaimed)
		if err := ldHandler.UpdateStatus(); err != nil {
			r.Recorder.Eventf(&ldHandler.ld, v1.EventTypeWarning, "UpdateStatusFail", "Update status fail, due to error: %v", err)
			log.WithError(err).Errorf("Update LocalDisk %v status fail", ldHandler.ld.Name)
			return reconcile.Result{}, err
		}
	}

	// At this stage, we have no relevant inspection data, so we won't do any processing for the time being
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *hwameistorv1alpha1.LocalDisk) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

// LocalDiskHandler
type LocalDiskHandler struct {
	client.Client
	record.EventRecorder
	ld     ldmv1alpha1.LocalDisk
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
func (ldHandler *LocalDiskHandler) GetLocalDisk(key client.ObjectKey) (*ldmv1alpha1.LocalDisk, error) {
	ld := ldmv1alpha1.LocalDisk{}
	if err := ldHandler.Get(context.Background(), key, &ld); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &ld, nil
}

// ListLocalDisk
func (ldHandler *LocalDiskHandler) ListLocalDisk() (*ldmv1alpha1.LocalDiskList, error) {
	list := &ldmv1alpha1.LocalDiskList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LocalDisk",
			APIVersion: "v1alpha1",
		},
	}

	err := ldHandler.List(context.TODO(), list)
	return list, err
}

// For
func (ldHandler *LocalDiskHandler) For(ld ldmv1alpha1.LocalDisk) *LocalDiskHandler {
	ldHandler.ld = ld
	ldHandler.filter = filter.NewLocalDiskFilter(ld)
	return ldHandler
}

// Bounded
func (ldHandler *LocalDiskHandler) UnClaimed() bool {
	return !ldHandler.filter.
		Init().
		Unclaimed().
		GetTotalResult()
}

// BoundTo assign disk to ldc
func (ldHandler *LocalDiskHandler) BoundTo(ldc ldmv1alpha1.LocalDiskClaim) error {
	ldcRef, err := reference.GetReference(nil, &ldc)
	if err != nil {
		return err
	}

	ldHandler.ld.Spec.ClaimRef = ldcRef
	ldHandler.ld.Status.State = ldmv1alpha1.LocalDiskClaimed

	if err = ldHandler.UpdateStatus(); err != nil {
		return err
	}
	ldHandler.EventRecorder.Eventf(&ldHandler.ld, v1.EventTypeNormal, "LocalDiskClaimed", "Claimed by %v/%v", ldc.Namespace, ldc.Name)
	return nil
}

// UpdateStatus
func (ldHandler *LocalDiskHandler) SetupStatus(status ldmv1alpha1.LocalDiskClaimState) {
	ldHandler.ld.Status.State = status
}

// UpdateStatus
func (ldHandler *LocalDiskHandler) UpdateStatus() error {
	return ldHandler.Update(context.Background(), &ldHandler.ld)
}

// ClaimRef
func (ldHandler *LocalDiskHandler) ClaimRef() *v1.ObjectReference {
	return ldHandler.ld.Spec.ClaimRef
}

// FilterDisk
func (ldHandler *LocalDiskHandler) FilterDisk(ldc ldmv1alpha1.LocalDiskClaim) bool {
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
