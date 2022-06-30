package node

import (
	"context"
	"fmt"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	ldmv1alpha1 "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	diskmonitor "github.com/hwameistor/local-storage/pkg/member/node/diskmonitor"
	isapisv1alpha1 "github.com/hwameistor/reliable-helper-system/pkg/apis/hwameistor/v1alpha1"
	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (m *manager) startLocalDiskTaskWorker(stopCh <-chan struct{}) {

	m.logger.Debug("LocalDisk Worker is working now")
	go func() {
		for {
			task, shutdown := m.localDiskTaskQueue.Get()
			if shutdown {
				m.logger.WithFields(log.Fields{"task": task}).Debug("Stop the LocalDisk worker")
				break
			}
			if err := m.processLocalDisk(task); err != nil {
				m.logger.WithFields(log.Fields{"task": task, "error": err.Error()}).Error("Failed to process LocalDisk task, retry later")
				m.localDiskTaskQueue.AddRateLimited(task)
			} else {
				m.logger.WithFields(log.Fields{"task": task}).Debug("Completed a LocalDisk task.")
				m.localDiskTaskQueue.Forget(task)
			}
			m.localDiskTaskQueue.Done(task)
		}
	}()

	<-stopCh
	m.volumeReplicaTaskQueue.Shutdown()
}

func (m *manager) processLocalDisk(localDiskNameSpacedName string) error {
	m.logger.Debug("processLocalDisk start ...")
	logCtx := m.logger.WithFields(log.Fields{"LocalDisk": localDiskNameSpacedName})
	logCtx.Debug("Working on a LocalDisk task")
	splitRes := strings.Split(localDiskNameSpacedName, "/")
	var diskName string
	if len(splitRes) >= 2 {
		// nameSpace = splitRes[0]
		diskName = splitRes[1]
	}

	localDisk := &ldmv1alpha1.LocalDisk{}
	if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Name: diskName}, localDisk); err != nil {
		if !errors.IsNotFound(err) {
			logCtx.WithError(err).Error("Failed to get LocalDisk from cache, retry it later ...")
			return err
		}
		logCtx.Info("Not found the LocalDisk from cache, should be deleted already.")
		return nil
	}

	m.logger.Debugf("Required node name %s, current node name %s.", localDisk.Spec.NodeName, m.name)
	if localDisk.Spec.NodeName != m.name {
		return nil
	}

	// check replacedisk label exists
	isInReplaceDiskTask, err := m.checkLocalDiskInReplaceDiskTaskByUUID(m.name, localDisk.Spec.UUID)
	//m.logger.Debugf("checkLocalDiskInReplaceDiskTaskByUUID isInReplaceDiskTask = %v, err = %v", isInReplaceDiskTask, err)
	if err != nil {
		//logCtx.Error("processLocalDisk checkLocalDiskInReplaceDiskTaskByUUID failed ,err = %v", err)
		return err
	}
	if !isInReplaceDiskTask {
		return nil
	}

	switch localDisk.Spec.State {
	case ldmv1alpha1.LocalDiskInactive:
		logCtx.Debug("LocalDiskInactive, todo ...")
		// 构建离线的event
		event := &diskmonitor.DiskEvent{}
		m.diskEventQueue.Add(event)
		return nil

	case ldmv1alpha1.LocalDiskActive:
		logCtx.Debug("LocalDiskActive ...")

	case ldmv1alpha1.LocalDiskUnknown:
		logCtx.Debug("LocalDiskUnknown ...")
		return nil

	default:
		logCtx.Error("Invalid LocalDisk state")
		return nil
	}

	switch localDisk.Status.State {
	case ldmv1alpha1.LocalDiskUnclaimed:
		//logCtx.Debug("LocalDiskUnclaimed localDisk = %v...", localDisk)
		err := m.processReplacedLocalDisk(localDisk)
		if err != nil {
			return err
		}
		key := client.ObjectKey{Name: localDisk.Name, Namespace: ""}
		ldisk, err := m.ldhandler.GetLocalDisk(key)
		if err != nil {
			return err
		}
		ldhandler := m.ldhandler.For(*ldisk)
		ldhandler.SetupStatus(ldmv1alpha1.LocalDiskClaimed)
		if err := m.ldhandler.UpdateStatus(); err != nil {
			log.WithError(err).Errorf("Update LocalDisk %v status fail", localDisk.Name)
			return err
		}
		return nil

	case ldmv1alpha1.LocalDiskReleased:
		//logCtx.Debug("LocalDiskReleased localDisk = %v...", localDisk)
		err := m.processReplacedLocalDisk(localDisk)
		if err != nil {
			return err
		}
		return nil

	case ldmv1alpha1.LocalDiskClaimed:
		logCtx.Debug("LocalDiskClaimed ...")
		return nil

	case ldmv1alpha1.LocalDiskInUse:
		logCtx.Debug("LocalDiskInUse ...")
		return nil

	case ldmv1alpha1.LocalDiskReserved:
		logCtx.Debug("LocalDiskReserved ...")
		return nil

	default:
		logCtx.Error("Invalid LocalDisk state")
	}

	return fmt.Errorf("invalid LocalDisk state")
}

func (m *manager) convertLdmLocalDisk(ld *ldmv1alpha1.LocalDisk) []*apisv1alpha1.LocalDisk {
	//m.logger.Debug("convertLocalDisk start ld = %v...", ld)
	if ld == nil {
		return nil
	}
	devicePath := ld.Spec.DevicePath
	if devicePath == "" || !strings.HasPrefix(devicePath, "/dev") || strings.Contains(devicePath, "mapper") {
		return nil
	}
	//if ld.Spec.State != ldmv1alpha1.LocalDiskActive {
	//	return nil
	//}

	var localDiskList = []*apisv1alpha1.LocalDisk{}
	disk := &apisv1alpha1.LocalDisk{}
	if ld.Spec.HasPartition {
		disk.State = apisv1alpha1.DiskStateInUse
	} else {
		disk.State = apisv1alpha1.DiskStateAvailable
	}
	disk.CapacityBytes = ld.Spec.Capacity
	disk.DevPath = devicePath
	disk.Class = ld.Spec.DiskAttributes.Type
	localDiskList = append(localDiskList, disk)
	return localDiskList
}

func (m *manager) processReplacedLocalDisk(ld *ldmv1alpha1.LocalDisk) error {
	//m.logger.Debug("processReplacedLocalDisk start ld = %+v...", ld)

	localDiskList := m.convertLdmLocalDisk(ld)
	if err := m.storageMgr.PoolManager().ExtendPools(localDiskList); err != nil {
		log.WithError(err).Error("Failed to ExtendPools")
		return err
	}

	// 如果lvm扩容失败，就不会执行如下同步资源流程
	localDisks, err := m.getLocalDisksMapByLdmLocalDisk(ld)
	if err != nil {
		log.WithError(err).Error("Failed to getLocalDisksMapByLocalDiskClaim")
		return err
	}
	//m.logger.Debug("processReplacedLocalDisk getLocalDisksMapByLdmLocalDisk  localDisks = %v, ld = %v", localDisks, ld)

	if err := m.storageMgr.Registry().SyncResourcesToNodeCRD(localDisks); err != nil {
		log.WithError(err).Error("Failed to SyncResourcesToNodeCRD")
		return err
	}
	return nil
}

func (m *manager) getLocalDisksMapByLdmLocalDisk(ld *ldmv1alpha1.LocalDisk) (map[string]*apisv1alpha1.LocalDisk, error) {
	m.logger.Debug("getLocalDisksMapByLocalDiskClaim ...")
	disks := make(map[string]*apisv1alpha1.LocalDisk)
	localDiskList := m.convertLdmLocalDisk(ld)
	for _, localDisk := range localDiskList {
		if localDisk == nil {
			continue
		}
		disks[localDisk.DevPath] = localDisk
	}
	return disks, nil
}

func (m *manager) checkLocalDiskInReplaceDiskTaskByUUID(nodeName, diskUuid string) (bool, error) {
	//m.logger.Debug("checkLocalDiskInReplaceDiskTaskByUUID nodeName = %v, diskUuid = %v", nodeName, diskUuid)
	replaceDiskList, err := m.rdhandler.ListReplaceDisk()
	if err != nil {
		m.logger.WithError(err).Error("Failed to get ReplaceDiskList")
		return false, err
	}
	//m.logger.Debug("checkLocalDiskInReplaceDiskTaskByUUID replaceDiskList %v", replaceDiskList)
	for _, replacedisk := range replaceDiskList.Items {
		if replacedisk.Spec.NodeName == nodeName {
			if replacedisk.Spec.ReplaceDiskStage != isapisv1alpha1.ReplaceDiskStage_Succeed && replacedisk.Spec.ReplaceDiskStage != isapisv1alpha1.ReplaceDiskStage_Failed {
				oldDiskUuid := replacedisk.Spec.OldUUID
				newDiskUuid := replacedisk.Spec.NewUUID
				//m.logger.Debug("checkLocalDiskInReplaceDiskTaskByUUID oldDiskUuid = %v,==== replacename = %v, ====newDiskUuid = %v, diskUuid = %v", oldDiskUuid, replacedisk.Name, newDiskUuid, diskUuid)
				if diskUuid == oldDiskUuid || diskUuid == newDiskUuid {
					return true, nil
				}
			}
		}
	}
	m.logger.Debug("checkLocalDiskInReplaceDiskTaskByUUID false")
	return false, nil
}
