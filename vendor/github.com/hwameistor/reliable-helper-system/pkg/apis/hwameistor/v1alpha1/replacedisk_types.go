package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReplaceDiskSpec defines the desired state of ReplaceDisk
type ReplaceDiskSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// UUID global unique identifier of the questioned disk
	OldUUID string `json:"oldUuid,omitempty"`
	// NewUUID global unique identifier of the new replaced disk
	NewUUID string `json:"newUuid,omitempty"`
	// NodeName nodeName of the replaced disk
	NodeName    string `json:"nodeName,omitempty"`
	DriverGroup string `json:"driverGroup,omitempty"`
	CtrId       string `json:"ctrId,omitempty"`
	SltId       string `json:"sltId,omitempty"`
	EID         string `json:"eId,omitempty"`
	// Init WaitDiskReplaced WaitSvcRestor Succeed
	ReplaceDiskStage ReplaceDiskStage `json:"replaceDiskStage,omitempty"`
	// IgnoreUnconvertibleVolumeData
	IgnoreUnconvertibleVolumeData bool `json:"ignoreUnconvertibleVolumeData,omitempty"`
}

// ReplaceDiskStage defines the observed state of replaceDisk
type ReplaceDiskStage string

const (
	// ReplaceDisk_Init represents that the disk at present locates Init stage.
	ReplaceDiskStage_Init ReplaceDiskStage = "Init"

	// ReplaceDisk_WaitDiskReplaced represents that the disk at present locates WaitDiskReplaced stage.
	ReplaceDiskStage_WaitDiskReplaced ReplaceDiskStage = "WaitDiskReplaced"

	// ReplaceDisk_WaitSvcRestor represents that the disk at present locates WaitSvcRestor stage.
	ReplaceDiskStage_WaitSvcRestor ReplaceDiskStage = "WaitSvcRestor"

	// ReplaceDisk_Succeed represents that the disk at present locates Succeed stage.
	ReplaceDiskStage_Succeed ReplaceDiskStage = "Succeed"

	// ReplaceDisk_Succeed represents that the disk at present locates Failed stage.
	ReplaceDiskStage_Failed ReplaceDiskStage = "Failed"
)

func (r ReplaceDiskStage) String() string {
	switch r {
	case ReplaceDiskStage_Init:
		return "Init"
	case ReplaceDiskStage_WaitDiskReplaced:
		return "WaitDiskReplaced"
	case ReplaceDiskStage_WaitSvcRestor:
		return "WaitSvcRestor"
	case ReplaceDiskStage_Succeed:
		return "Succeed"
	case ReplaceDiskStage_Failed:
		return "Failed"
	default:
		return "NA"
	}
}

// ReplaceDiskStatus defines the observed status of OldReplaceDisk and NewReplaceDisk
type ReplaceDiskStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Init WaitDataRepair WaitDiskLVMRelease WaitDiskRemoved
	OldDiskReplaceStatus ReplaceStatus `json:"oldDiskReplaceStatus,omitempty"`
	// Init WaitDiskLVMRejoin WaitDataBackup Succeed
	NewDiskReplaceStatus ReplaceStatus `json:"newDiskReplaceStatus,omitempty"`
	// MigrateVolumeNames migrateVolumeNames of the replaced disk
	MigrateVolumeNames []string `json:"migrateVolumeNames,omitempty"`
	// MigrateSucceededVolumeNames migrateSucceededVolumeNames of the replaced disk
	MigrateSucceededVolumeNames []string `json:"migrateSucceededVolumeNames,omitempty"`
	// MigrateBackUpSucceededVolumeNames migrateBackUpSucceededVolumeNames of the replaced disk
	MigrateBackUpSucceededVolumeNames []string `json:"migrateBackUpSucceededVolumeNames,omitempty"`
	// MigrateFailededVolumeNames migrateFailededVolumeNames of the replaced disk
	MigrateFailededVolumeNames []string `json:"migrateFailededVolumeNames,omitempty"`
	// MigrateBackUpFailededVolumeNames migrateBackUpFailededVolumeNames of the replaced disk
	MigrateBackUpFailededVolumeNames []string `json:"migrateBackUpFailededVolumeNames,omitempty"`
	// ErrMsg errMsg of the replaced disk
	ErrMsg string `json:"errMsg,omitempty"`
	// WarnMsg warnMsg of the replaced disk
	WarnMsg string `json:"warnMsg,omitempty"`
}

// ReplaceStatus defines the observed status of replacedDisk
type ReplaceStatus string

const (
	// ReplaceDisk_Init represents that the disk at present locates Init status.
	ReplaceDisk_Init ReplaceStatus = "Init"

	// ReplaceDisk_WaitDataRepair represents that the disk at present locates WaitDataRepair status.
	ReplaceDisk_WaitDataRepair ReplaceStatus = "WaitDataRepair"

	// ReplaceDisk_WaitDiskLVMRelease represents that the disk at present locates WaitDiskLVMRelease status.
	ReplaceDisk_WaitDiskLVMRelease ReplaceStatus = "WaitDiskLVMRelease"

	// ReplaceDisk_DiskLVMReleased represents that the disk at present locates DiskLVMReleased status.
	ReplaceDisk_DiskLVMReleased ReplaceStatus = "DiskLVMReleased"

	// ReplaceDisk_WaitDiskLVMRejoin represents that the disk at present locates WaitDiskLVMRejoin status.
	ReplaceDisk_WaitDiskLVMRejoin ReplaceStatus = "WaitDiskLVMRejoin"

	// ReplaceDisk_WaitDataBackup represents that the disk at present locates WaitDataBackup status.
	ReplaceDisk_WaitDataBackup ReplaceStatus = "WaitDataBackup"

	// ReplaceDisk_DataBackuped represents that the disk at present locates DataBackuped status.
	ReplaceDisk_DataBackuped ReplaceStatus = "DataBackuped"

	// ReplaceDisk_Succeed represents that the disk at present locates Succeed status.
	ReplaceDisk_Succeed ReplaceStatus = "Succeed"

	// ReplaceDisk_Failed represents that the disk at present locates Failed status.
	ReplaceDisk_Failed ReplaceStatus = "Failed"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReplaceDisk is the Schema for the replacedisks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=replacedisks,scope=Namespaced
type ReplaceDisk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplaceDiskSpec   `json:"spec,omitempty"`
	Status ReplaceDiskStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReplaceDiskList contains a list of ReplaceDisk
type ReplaceDiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplaceDisk `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplaceDisk{}, &ReplaceDiskList{})
}
