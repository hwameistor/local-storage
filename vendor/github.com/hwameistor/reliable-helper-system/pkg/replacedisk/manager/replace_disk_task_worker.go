package manager

import (
	"context"
	errs "errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	ldm "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-disk-manager/pkg/controller/localdisk"
	lsapisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	apisv1alpha1 "github.com/hwameistor/reliable-helper-system/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/reliable-helper-system/pkg/utils"
	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (m *manager) startReplaceDiskTaskWorker(stopCh <-chan struct{}) {
	m.logger.Debug("ReplaceDisk Worker is working now")
	go func() {
		for {
			task, shutdown := m.ReplaceDiskNodeManager().ReplaceDiskTaskQueue().Get()
			m.logger.Debugf("startReplaceDiskTaskWorker task = %+v", task)
			if shutdown {
				m.logger.WithFields(log.Fields{"task": task}).Debug("Stop the ReplaceDisk worker")
				break
			}
			if err := m.processReplaceDisk(task); err != nil {
				m.logger.WithFields(log.Fields{"task": task, "error": err.Error()}).Error("Failed to process ReplaceDisk task, retry later")
				m.ReplaceDiskNodeManager().ReplaceDiskTaskQueue().AddRateLimited(task)
			} else {
				m.logger.WithFields(log.Fields{"task": task}).Debug("Completed a ReplaceDisk task.")
				m.ReplaceDiskNodeManager().ReplaceDiskTaskQueue().Forget(task)
			}
			m.ReplaceDiskNodeManager().ReplaceDiskTaskQueue().Done(task)
		}
	}()

	<-stopCh
	m.ReplaceDiskNodeManager().ReplaceDiskTaskQueue().Shutdown()
}

func (m *manager) processReplaceDisk(replaceDiskNameSpacedName string) error {
	m.logger.Debug("processReplaceDisk start ...")
	logCtx := m.logger.WithFields(log.Fields{"ReplaceDisk": replaceDiskNameSpacedName})

	logCtx.Debug("Working on a ReplaceDisk task")
	splitRes := strings.Split(replaceDiskNameSpacedName, "/")
	var nameSpace, replaceDiskName string
	if len(splitRes) >= 2 {
		nameSpace = splitRes[0]
		replaceDiskName = splitRes[1]
	}
	replaceDisk := &apisv1alpha1.ReplaceDisk{}
	if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Namespace: nameSpace, Name: replaceDiskName}, replaceDisk); err != nil {
		if !errors.IsNotFound(err) {
			logCtx.WithError(err).Error("Failed to get ReplaceDisk from cache, retry it later ...")
			return err
		}
		//logCtx.Info("Not found the ReplaceDisk from cache, should be deleted already. err = %v", err)
		fmt.Printf("Not found the ReplaceDisk from cache, should be deleted already. err = %v", err)
		return nil
	}

	rdhandler := m.rdhandler.SetReplaceDisk(*replaceDisk)
	m.rdhandler = rdhandler

	m.logger.Debugf("Required node name %s, current node name %s.", replaceDisk.Spec.NodeName, m.nodeName)
	if replaceDisk.Spec.NodeName != m.nodeName {
		return nil
	}

	m.logger.Debugf("processReplaceDisk replaceDisk.Status %v.", replaceDisk.Status)
	switch replaceDisk.Status.OldDiskReplaceStatus {
	case apisv1alpha1.ReplaceDisk_Init:
		// do init job
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_Init && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_Init {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_WaitDataRepair
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_Init
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_WaitDataRepair,ReplaceDisk_Init failed")
				return err
			}
		}
		// update init stage
		err := m.processOldReplaceDiskStatusInit(replaceDisk)
		if err != nil {
			m.logger.Errorf("processOldReplaceDiskStatusInit failed err = %v", err)
			m.rdhandler.SetErrMsg(err.Error())
			return err
		}
		return nil
	case apisv1alpha1.ReplaceDisk_WaitDataRepair:
		err := m.processOldReplaceDiskStatusWaitDataRepair(replaceDisk)
		if err != nil {
			log.Error(err, "processReplaceDisk processOldReplaceDiskStatusWaitDataRepair failed")
			return err
		}
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_WaitDataRepair && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_Init {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_WaitDiskLVMRelease
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_Init
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				m.rdhandler.SetErrMsg(err.Error())
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_WaitDiskLVMRelease,ReplaceDisk_Init failed")
				return err
			}
		}

		return nil
	case apisv1alpha1.ReplaceDisk_WaitDiskLVMRelease:
		err := m.processOldReplaceDiskStatusWaitDiskLVMRelease(replaceDisk)
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if err != nil {
			log.Error(err, "processReplaceDisk processOldReplaceDiskStatusWaitDiskLVMRelease failed")
			return err
		}
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_WaitDiskLVMRelease && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_Init {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DiskLVMReleased
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_Init
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_DiskLVMReleased,ReplaceDisk_Init failed")
				return err
			}
		}
		return nil
	case apisv1alpha1.ReplaceDisk_DiskLVMReleased:
		if replaceDisk.Spec.ReplaceDiskStage == apisv1alpha1.ReplaceDiskStage_Init || replaceDisk.Spec.ReplaceDiskStage == apisv1alpha1.ReplaceDiskStage_WaitDiskReplaced {
			return m.processOldReplaceDiskStatusDiskLVMReleased(replaceDisk)
		}
	case apisv1alpha1.ReplaceDisk_Failed:
		return m.processOldReplaceDiskStatusFailed(replaceDisk)
	default:
		logCtx.Error("Invalid ReplaceDisk status")
		return errs.New("invalid ReplaceDisk status")
	}

	if replaceDisk.Status.OldDiskReplaceStatus != apisv1alpha1.ReplaceDisk_DiskLVMReleased {
		logCtx.Errorf("Invalid ReplaceDisk OldDiskReplaceStatus,replaceDisk.Status.OldDiskReplaceStatus = %v", replaceDisk.Status.OldDiskReplaceStatus)
		return errs.New("invalid ReplaceDisk OldDiskReplaceStatus")
	}

	switch replaceDisk.Status.NewDiskReplaceStatus {
	case apisv1alpha1.ReplaceDisk_Init:
		// todo create job formating disk
		// 触发ld更新
		err := m.updateLocalDisk(replaceDisk)
		if err != nil {
			log.Error(err, "processReplaceDisk NewDiskReplaceStatus updateLocalDisk failed")
			return err
		}
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_DiskLVMReleased && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_Init {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DiskLVMReleased
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_WaitDiskLVMRejoin
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "ReplaceDiskStage_WaitSvcRestor UpdateReplaceDiskStatus ReplaceDisk_DiskLVMReleased,ReplaceDisk_WaitDiskLVMRejoin failed")
				return err
			}
		}
		return nil
	case apisv1alpha1.ReplaceDisk_WaitDiskLVMRejoin:
		err := m.processNewReplaceDiskStatusWaitDiskLVMRejoin(replaceDisk)
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if err != nil {
			log.Error(err, "processReplaceDisk processNewReplaceDiskStatusWaitDiskLVMRejoin failed")
			m.rdhandler.SetErrMsg(err.Error())
			return err
		}
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_DiskLVMReleased && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_WaitDiskLVMRejoin {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DiskLVMReleased
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_WaitDataBackup
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_DiskLVMReleased,ReplaceDisk_WaitDataBackup failed")
				return err
			}
		}

		return nil
	case apisv1alpha1.ReplaceDisk_WaitDataBackup:
		err := m.processNewReplaceDiskStatusWaitDataBackup(replaceDisk)
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if err != nil {
			log.Error(err, "processReplaceDisk processNewReplaceDiskStatusWaitDataBackup failed")
			m.rdhandler.SetErrMsg(err.Error())
			return err
		}
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_DiskLVMReleased && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_WaitDataBackup {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DiskLVMReleased
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DataBackuped
			if err := rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_DiskLVMReleased,ReplaceDisk_DataBackuped failed")
				return err
			}
		}

		return nil
	case apisv1alpha1.ReplaceDisk_DataBackuped:
		replaceDiskStatus := rdhandler.ReplaceDiskStatus()
		if replaceDiskStatus.OldDiskReplaceStatus == apisv1alpha1.ReplaceDisk_DiskLVMReleased && replaceDiskStatus.NewDiskReplaceStatus == apisv1alpha1.ReplaceDisk_DataBackuped {
			replaceDiskStatus.OldDiskReplaceStatus = apisv1alpha1.ReplaceDisk_DiskLVMReleased
			replaceDiskStatus.NewDiskReplaceStatus = apisv1alpha1.ReplaceDisk_Succeed
			if err := m.rdhandler.UpdateReplaceDiskStatus(replaceDiskStatus); err != nil {
				log.Error(err, "processReplaceDisk UpdateReplaceDiskStatus ReplaceDisk_DiskLVMReleased,ReplaceDisk_Succeed failed")
				return err
			}
		}
		err := m.processNewReplaceDiskStatusDataBackuped(replaceDisk)
		if err != nil {
			log.Error(err, "processReplaceDisk processNewReplaceDiskStatusDataBackuped failed")
			m.rdhandler.SetErrMsg(err.Error())
			return err
		}
		return nil
	case apisv1alpha1.ReplaceDisk_Succeed:
		return m.processNewReplaceDiskStatusSucceed(replaceDisk)
	case apisv1alpha1.ReplaceDisk_Failed:
		return m.processNewReplaceDiskStatusFailed(replaceDisk)
	default:
		logCtx.Error("Invalid ReplaceDisk status")
	}

	return fmt.Errorf("invalid ReplaceDisk status")
}

func (m *manager) processOldReplaceDiskStatusInit(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	replaceDisk.Spec.ReplaceDiskStage = apisv1alpha1.ReplaceDiskStage_WaitDiskReplaced
	m.rdhandler.SetReplaceDiskStage(replaceDisk.Spec.ReplaceDiskStage)
	return m.rdhandler.UpdateReplaceDiskCR()
}

func (m *manager) processOldReplaceDiskStatusWaitDataRepair(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("processOldReplaceDiskStatusWaitDataRepair start ... ")

	logCtx := m.logger.WithFields(log.Fields{"ReplaceDisk": replaceDisk.Name})

	oldDiskName, err := m.getDiskNameByDiskUUID(replaceDisk.Spec.OldUUID, replaceDisk.Spec.NodeName)
	if err != nil {
		logCtx.Errorf("processOldReplaceDiskStatusWaitDataRepair getDiskNameByDiskUUID failed err = %v", err)
		return err
	}

	oldLocalDisk, err := m.getLocalDiskByDiskName(oldDiskName, replaceDisk.Spec.NodeName)
	if err != nil {
		m.logger.WithError(err).Error("processOldReplaceDiskStatusWaitDataRepair: Failed to getLocalDiskByDiskName")
		return err
	}

	// directly replacedisk ; not do datarepair
	if oldLocalDisk.Spec.HasRAID == true {
		m.logger.Debug("processOldReplaceDiskStatusWaitDataRepair oldLocalDisk.Spec.HasRAID is true.")
		return nil
	}

	localVolumeReplicasMap, err := m.getAllLocalVolumeAndReplicasMapOnDisk(oldDiskName, replaceDisk.Spec.NodeName)
	if err != nil {
		logCtx.Errorf("processOldReplaceDiskStatusWaitDataRepair getAllLocalVolumeAndReplicasMapOnDisk failed err = %v", err)
		return err
	}

	var migrateVolumeNames []string
	var localVolumeList []lsapisv1alpha1.LocalVolume
	for localVolumeName := range localVolumeReplicasMap {
		vol := &lsapisv1alpha1.LocalVolume{}
		if err = m.apiClient.Get(context.TODO(), types.NamespacedName{Name: localVolumeName}, vol); err != nil {
			if !errors.IsNotFound(err) {
				log.Errorf("Failed to get Volume from cache, retry it later ...")
				return err
			}
			log.Printf("Not found the Volume from cache, should be deleted already.")
			return nil
		}

		if vol.Spec.ReplicaNumber == 1 && vol.Spec.Convertible == false {
			err := fmt.Sprintf("The volume %v replicanumber is 1 and convertible is false , cannot do replacedisk task for the risk of dataloss", vol.Name)
			m.rdhandler.SetErrMsg(err)
			continue
		}

		err = m.createMigrateTaskByLocalVolume(*vol)
		if err != nil {
			log.Error(err, "createMigrateTaskByLocalVolume failed")
			return err
		}
		migrateVolumeNames = append(migrateVolumeNames, vol.Name)

		m.rdhandler.SetMigrateVolumeNames(migrateVolumeNames)
		localVolumeList = append(localVolumeList, *vol)
	}

	if len(localVolumeList) == 0 {
		return nil
	}

	err = m.waitMigrateTaskByLocalVolumeDone(localVolumeList, MigrateType_DataRepair)
	m.logger.Debug("processOldReplaceDiskStatusWaitDataRepair end ... ")
	return err
}

func (m *manager) createMigrateTaskByLocalVolume(vol lsapisv1alpha1.LocalVolume) error {
	m.logger.Debug("createMigrateTaskByLocalVolume start ... ")

	localVolumeMigrate, err := m.migrateCtr.ConstructLocalVolumeMigrate(vol)
	if err != nil {
		log.Error(err, "createMigrateTaskByLocalVolume ConstructLocalVolumeMigrate failed")
		return err
	}
	migrate := &lsapisv1alpha1.LocalVolumeMigrate{}
	if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Name: localVolumeMigrate.Name}, migrate); err == nil {
		m.logger.Info("createMigrateTaskByLocalVolume localVolumeMigrate %v already exists.", localVolumeMigrate.Name)
		return nil
	}

	err = m.migrateCtr.CreateLocalVolumeMigrate(*localVolumeMigrate)
	if err != nil {
		log.Error(err, "createMigrateTaskByLocalVolume CreateLocalVolumeMigrate failed")
		return err
	}

	m.logger.Debug("createMigrateTaskByLocalVolume end ... ")
	return nil
}

// 统计数据迁移情况
func (m *manager) waitMigrateTaskByLocalVolumeDone(volList []lsapisv1alpha1.LocalVolume, migrateType string) error {
	m.logger.Debug("waitMigrateTaskByLocalVolumeDone start ... ")
	var wg sync.WaitGroup
	var migrateSucceedVols = m.rdhandler.ReplaceDiskStatus().MigrateSucceededVolumeNames
	var migrateBackUpSucceedVols = m.rdhandler.ReplaceDiskStatus().MigrateBackUpSucceededVolumeNames
	var migrateFailedVols = m.rdhandler.ReplaceDiskStatus().MigrateFailededVolumeNames
	var migrateBackUpFailedVols = m.rdhandler.ReplaceDiskStatus().MigrateBackUpFailededVolumeNames
	for _, vol := range volList {
		vol := vol
		wg.Add(1)
		go func() {
			for {
				migrateStatus, err := m.migrateCtr.GetLocalVolumeMigrateStatusByLocalVolume(vol)
				if err != nil {
					log.Error(err, "waitMigrateTaskByLocalVolumeDone GetLocalVolumeMigrateStatusByLocalVolume failed")
					return
				}
				if migrateStatus.State == lsapisv1alpha1.OperationStateSubmitted || migrateStatus.State == lsapisv1alpha1.OperationStateInProgress {
					time.Sleep(2 * time.Second)
					continue
				}
				if migrateStatus.State == lsapisv1alpha1.OperationStateAborted || migrateStatus.State == lsapisv1alpha1.OperationStateAborting ||
					migrateStatus.State == lsapisv1alpha1.OperationStateToBeAborted {
					if migrateType == MigrateType_DataRepair {
						migrateFailedVols = append(migrateFailedVols, vol.Name)
					}
					if migrateType == MigrateType_BackUp {
						migrateBackUpFailedVols = append(migrateBackUpFailedVols, vol.Name)
					}
					continue
				}
				if migrateType == MigrateType_DataRepair {
					migrateSucceedVols = append(migrateSucceedVols, vol.Name)
				}
				if migrateType == MigrateType_BackUp {
					migrateBackUpSucceedVols = append(migrateBackUpSucceedVols, vol.Name)
				}
				break
			}
			wg.Done()
		}()
	}
	wg.Wait()

	err := m.rdhandler.SetMigrateSucceededVolumeNames(migrateSucceedVols).SetMigrateFailededVolumeNames(migrateFailedVols).
		SetMigrateBackUpSucceededVolumeNames(migrateBackUpSucceedVols).
		SetMigrateBackUpFailededVolumeNames(migrateBackUpFailedVols).
		UpdateReplaceDiskStatus(m.rdhandler.ReplaceDiskStatus())
	if err != nil {
		m.logger.Errorf("waitMigrateTaskByLocalVolumeDone UpdateReplaceDiskStatus failed ... ")
		return err
	}

	m.logger.Debug("waitMigrateTaskByLocalVolumeDone end ... ")
	return nil
}

func (m *manager) getLocalDiskByDiskName(diskName, nodeName string) (ldm.LocalDisk, error) {
	m.logger.Debug("getLocalDiskByDiskName start ... ")
	// replacedDiskName e.g.(/dev/sdb -> sdb)
	var replacedDiskName string
	if strings.HasPrefix(diskName, "/dev") {
		replacedDiskName = strings.Replace(diskName, "/dev/", "", 1)
	}

	// ConvertNodeName e.g.(10.23.10.12 => 10-23-10-12)
	localDiskName := utils.ConvertNodeName(nodeName) + "-" + replacedDiskName
	key := client.ObjectKey{Name: localDiskName, Namespace: ""}
	localDisk, err := m.localDiskController.GetLocalDisk(key)
	if err != nil {
		m.logger.WithError(err).Error("getLocalDiskByDiskName: Failed to GetLocalDisk")
		return localDisk, err
	}
	m.logger.Debug("getLocalDiskByDiskName end ... ")
	return localDisk, nil
}

func (m *manager) processOldReplaceDiskStatusWaitDiskLVMRelease(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("processOldReplaceDiskStatusWaitDiskLVMRelease start ... ")
	diskUuid := replaceDisk.Spec.OldUUID
	nodeName := replaceDisk.Spec.NodeName
	diskName, err := m.getDiskNameByDiskUUID(diskUuid, nodeName)
	if err != nil {
		m.logger.WithError(err).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to getDiskNameByDiskUUID")
		return err
	}

	localDisk, err := m.getLocalDiskByDiskName(diskName, nodeName)
	if err != nil {
		m.logger.WithError(err).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to getLocalDiskByDiskName")
		return err
	}

	diskType := localDisk.Spec.DiskAttributes.Type
	volGroupName, err := utils.GetPoolNameAccordingDiskType(diskType)
	if err != nil {
		m.logger.WithError(err).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to GetPoolNameAccordingDiskType")
		return err
	}

	var options []string
	ctx := context.TODO()
	nodeList := &lsapisv1alpha1.LocalStorageNodeList{}
	if err := m.apiClient.List(ctx, nodeList); err != nil {
		m.logger.WithError(err).Error("Failed to get NodeList")
		return err
	}

	for _, node := range nodeList.Items {
		if node.Name == nodeName {
			disks := node.Status.Pools[volGroupName].Disks
			if len(disks) == 1 {
				errMsg := "VgGroup has only one disk, cannot do replacedisk operation"
				return errors.NewBadRequest(errMsg)
			}
			break
		}
	}

	err = m.cmdExec.vgreduce(volGroupName, diskName, options)
	if err != nil {
		m.logger.WithError(err).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to vgreduce")
		err1 := m.cmdExec.pvremove(volGroupName, diskName, options)
		if err1 != nil {
			m.logger.WithError(err1).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to pvremove")
			return err1
		}
		return err
	}

	// check if vgs can execute pvremove
	//err1 := m.cmdExec.pvremove(volGroupName, diskName, options)
	//if err1 != nil {
	//	m.logger.WithError(err1).Error("processOldReplaceDiskStatusWaitDiskLVMRelease: Failed to pvremove")
	//	return err1
	//}

	key := client.ObjectKey{Name: localDisk.Name, Namespace: ""}
	ldisk, err := m.ldhandler.GetLocalDisk(key)
	if err != nil {
		return err
	}

	ldhandler := m.ldhandler.For(*ldisk)
	ldhandler.SetupStatus(ldm.LocalDiskReleased)
	if err := m.ldhandler.UpdateStatus(); err != nil {
		log.WithError(err).Errorf("Update LocalDisk %v status fail", localDisk.Name)
		return err
	}

	m.logger.Debugf("processOldReplaceDiskStatusWaitDiskLVMRelease end ... ")
	return nil
}

func (m manager) processOldReplaceDiskStatusDiskLVMReleased(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	replaceDisk.Spec.ReplaceDiskStage = apisv1alpha1.ReplaceDiskStage_WaitSvcRestor
	m.rdhandler.SetReplaceDiskStage(replaceDisk.Spec.ReplaceDiskStage)
	return m.rdhandler.UpdateReplaceDiskCR()
}

func (m *manager) processOldReplaceDiskStatusFailed(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	return nil
}

func (m *manager) processNewReplaceDiskStatusWaitDiskLVMRejoin(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("processNewReplaceDiskStatusWaitDiskLVMRejoin start ... ")
	err := m.waitDiskLVMRejoinDone(replaceDisk)
	if err != nil {
		m.logger.WithError(err).Error("processNewReplaceDiskStatusWaitDiskLVMRejoin: Failed to waitDiskLVMRejoinDone")
		return err
	}
	m.logger.Debug("processNewReplaceDiskStatusWaitDiskLVMRejoin end ... ")

	return nil
}

func (m *manager) waitDiskLVMRejoinDone(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("waitDiskLVMRejoinDone start ... ")

	diskUuid := replaceDisk.Spec.NewUUID
	nodeName := replaceDisk.Spec.NodeName
	diskName, err := m.getDiskNameByDiskUUID(diskUuid, nodeName)
	if err != nil {
		m.logger.WithError(err).Error("waitDiskLVMRejoinDone: Failed to getDiskNameByDiskUUID")
		return err
	}

	var retryNums int

	for {
		retryNums++
		if m.rdhandler.CheckReplaceDiskTaskDestroyed(*replaceDisk) {
			break
		}
		localDisk, err := m.getLocalDiskByDiskName(diskName, nodeName)
		if err != nil {
			m.logger.WithError(err).Error("waitDiskLVMRejoinDone: Failed to getLocalDiskByDiskName")
			return err
		}

		m.logger.WithError(err).Errorf("waitDiskLVMRejoinDone getLocalDiskByDiskName localDisk.Status.State = %v, retryNums = %v", localDisk.Status.State, retryNums)
		if localDisk.Status.State == ldm.LocalDiskClaimed {
			break
		}

		if retryNums > 10 {
			// 触发ld更新
			err := m.updateLocalDisk(replaceDisk)
			if err != nil {
				log.Error(err, "processReplaceDisk NewDiskReplaceStatus updateLocalDisk failed")
				return err
			}
			retryNums = 0
		}
	}
	m.logger.Debug("waitDiskLVMRejoinDone end ... ")

	return nil
}

func (m *manager) processNewReplaceDiskStatusWaitDataBackup(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("processNewReplaceDiskStatusWaitDataBackup start ... ")

	if replaceDisk == nil {
		return errs.New("processNewReplaceDiskStatusWaitDataBackup replaceDisk is nil")
	}

	oldDiskName, err := m.getDiskNameByDiskUUID(replaceDisk.Spec.OldUUID, replaceDisk.Spec.NodeName)
	if err != nil {
		m.logger.Errorf("processNewReplaceDiskStatusWaitDataBackup getDiskNameByDiskUUID failed err = %v", err)
		return err
	}

	oldLocalDisk, err := m.getLocalDiskByDiskName(oldDiskName, replaceDisk.Spec.NodeName)
	if err != nil {
		m.logger.WithError(err).Error("processNewReplaceDiskStatusWaitDataBackup: Failed to getLocalDiskByDiskName")
		return err
	}

	// directly replacedisk ; no data backup
	if oldLocalDisk.Spec.HasRAID == true {
		m.logger.Debug("processNewReplaceDiskStatusWaitDataBackup oldLocalDisk.Spec.HasRAID is true.")
		return nil
	}

	migrateLocalVolumes := replaceDisk.Status.MigrateVolumeNames
	var localVolumeList []lsapisv1alpha1.LocalVolume
	m.logger.Debug("processNewReplaceDiskStatusWaitDataBackup migrateLocalVolumes = %v", migrateLocalVolumes)
	for _, localVolumeName := range migrateLocalVolumes {
		vol := &lsapisv1alpha1.LocalVolume{}
		if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Name: localVolumeName}, vol); err != nil {
			if !errors.IsNotFound(err) {
				log.Errorf("Failed to get Volume from cache, retry it later ...")
				return err
			}
			log.Printf("Not found the Volume from cache, should be deleted already.")
			return nil
		}
		err := m.createMigrateTaskByLocalVolume(*vol)
		if err != nil {
			log.Error(err, "processNewReplaceDiskStatusWaitDataBackup: createMigrateTaskByLocalVolume failed")
			return err
		}
		localVolumeList = append(localVolumeList, *vol)
	}

	err = m.waitMigrateTaskByLocalVolumeDone(localVolumeList, MigrateType_BackUp)
	if err != nil {
		m.logger.WithError(err).Error("processNewReplaceDiskStatusWaitDataBackup: Failed to waitMigrateTaskByLocalVolumeDone")
		return err
	}
	m.logger.Debug("processNewReplaceDiskStatusWaitDataBackup end ... ")

	return nil
}

func (m manager) processNewReplaceDiskStatusDataBackuped(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	replaceDisk.Spec.ReplaceDiskStage = apisv1alpha1.ReplaceDiskStage_Succeed
	m.rdhandler.SetReplaceDiskStage(replaceDisk.Spec.ReplaceDiskStage)
	return m.rdhandler.UpdateReplaceDiskCR()
}

// 统计结果并上报
func (m *manager) processNewReplaceDiskStatusSucceed(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	if replaceDisk == nil {
		return errs.New("processNewReplaceDiskStatusSucceed replaceDisk is nil")
	}

	oldDiskName, err := m.getDiskNameByDiskUUID(replaceDisk.Spec.OldUUID, replaceDisk.Spec.NodeName)
	if err != nil {
		m.logger.Errorf("processNewReplaceDiskStatusSucceed getDiskNameByDiskUUID failed err = %v", err)
		return err
	}

	oldLocalDisk, err := m.getLocalDiskByDiskName(oldDiskName, replaceDisk.Spec.NodeName)
	if err != nil {
		m.logger.WithError(err).Error("processNewReplaceDiskStatusSucceed: Failed to getLocalDiskByDiskName")
		return err
	}

	// directly replacedisk ; no data backup
	if oldLocalDisk.Spec.HasRAID == true {
		m.logger.Debug("processNewReplaceDiskStatusSucceed oldLocalDisk.Spec.HasRAID is true.")
		err := fmt.Sprintf("Please Check Raid Device Status, If Not Ready, Repair It !")
		m.rdhandler.SetErrMsg(err)
		return nil
	}
	return nil
}

// 告警
func (m *manager) processNewReplaceDiskStatusFailed(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	return nil
}

func (m *manager) getAllLocalVolumeAndReplicasMapOnDisk(diskName, nodeName string) (map[string][]*lsapisv1alpha1.LocalVolumeReplica, error) {
	m.logger.Debug("getAllLocalVolumeAndReplicasMapOnDisk start ... ")

	replicaList := &lsapisv1alpha1.LocalVolumeReplicaList{}
	if err := m.apiClient.List(context.TODO(), replicaList); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	localVolumeReplicasMap := make(map[string][]*lsapisv1alpha1.LocalVolumeReplica)
	for _, replica := range replicaList.Items {
		if replica.Spec.NodeName != nodeName {
			continue
		}
		for _, dName := range replica.Status.Disks {
			// start with /dev/sdx
			if dName == diskName {
				replicas, ok := localVolumeReplicasMap[replica.Spec.VolumeName]
				if !ok {
					var replicaList []*lsapisv1alpha1.LocalVolumeReplica
					replicaList = append(replicaList, &replica)
					localVolumeReplicasMap[replica.Spec.VolumeName] = replicaList
				} else {
					replicas = append(replicas, &replica)
					localVolumeReplicasMap[replica.Spec.VolumeName] = replicas
				}
			}
		}
	}
	m.logger.Debug("getAllLocalVolumeAndReplicasMapOnDisk end ... ")

	return localVolumeReplicasMap, nil
}

func (m *manager) CheckVolumeInReplaceDiskTask(nodeName, volName string) (bool, error) {
	m.logger.Debug("CheckVolumeInReplaceDiskTask start ... ")

	allLocalVolumeReplicasMap, err := m.getAllLocalVolumesAndReplicasInRepalceDiskTask(nodeName)
	if err != nil {
		m.logger.WithError(err).Error("Failed to getAllLocalVolumesAndReplicasInRepalceDiskTask")
		return false, err
	}
	if _, ok := allLocalVolumeReplicasMap[volName]; ok {
		return true, nil
	}
	m.logger.Debug("CheckVolumeInReplaceDiskTask end ... ")

	return false, nil
}

func (m *manager) getAllLocalVolumesAndReplicasInRepalceDiskTask(nodeName string) (map[string][]*lsapisv1alpha1.LocalVolumeReplica, error) {
	m.logger.Debug("getAllLocalVolumesAndReplicasInRepalceDiskTask start ... ")

	diskNames, err := m.getDiskNamesInRepalceDiskTask(nodeName)
	if err != nil {
		m.logger.WithError(err).Error("Failed to get ReplaceDiskList")
		return nil, err
	}

	var allLocalVolumeReplicasMap = make(map[string][]*lsapisv1alpha1.LocalVolumeReplica)
	for _, diskName := range diskNames {
		localVolumeReplicasMap, err := m.getAllLocalVolumeAndReplicasMapOnDisk(diskName, nodeName)
		if err != nil {
			m.logger.WithError(err).Error("Failed to getAllLocalVolumeAndReplicasMapOnDisk")
			return nil, err
		}
		allLocalVolumeReplicasMap = utils.MergeLocalVolumeReplicaMap(allLocalVolumeReplicasMap, localVolumeReplicasMap)
	}

	m.logger.Debug("getAllLocalVolumesAndReplicasInRepalceDiskTask end ... ")
	return allLocalVolumeReplicasMap, nil
}

func (m *manager) getDiskNamesInRepalceDiskTask(nodeName string) ([]string, error) {
	m.logger.Debug("getDiskNamesInRepalceDiskTask start ... ")

	replaceDiskList, err := m.rdhandler.ListReplaceDisk()
	if err != nil {
		m.logger.WithError(err).Error("Failed to get ReplaceDiskList")
		return nil, err
	}

	var diskNames []string
	for _, replacedisk := range replaceDiskList.Items {
		if replacedisk.Spec.NodeName == nodeName {
			diskUuid := replacedisk.Spec.OldUUID
			diskName, err := m.getDiskNameByDiskUUID(diskUuid, nodeName)
			if err != nil {
				m.logger.WithError(err).Error("Failed to getDiskNameByDiskUUID")
				return nil, err
			}
			diskNames = append(diskNames, diskName)
		}
	}
	m.logger.Debug("getDiskNamesInRepalceDiskTask end ... ")

	return diskNames, nil
}

func (m *manager) CheckLocalDiskInReplaceDiskTaskByUUID(nodeName, diskUuid string) (bool, error) {
	m.logger.Debug("CheckLocalDiskInReplaceDiskTaskByUUID start ... ")

	replaceDiskList, err := m.rdhandler.ListReplaceDisk()
	if err != nil {
		m.logger.WithError(err).Error("Failed to get ReplaceDiskList")
		return false, err
	}

	for _, replacedisk := range replaceDiskList.Items {
		if replacedisk.Spec.NodeName == nodeName {
			oldDiskUuid := replacedisk.Spec.OldUUID
			newDiskUuid := replacedisk.Spec.NewUUID
			if diskUuid == oldDiskUuid || diskUuid == newDiskUuid {
				return true, nil
			}
		}
	}

	m.logger.Debug("CheckLocalDiskInReplaceDiskTaskByUUID end ... ")
	return false, nil
}

// start with /dev/sdx
func (m *manager) getDiskNameByDiskUUID(diskUUID, nodeName string) (string, error) {
	var recorder record.EventRecorder
	ldHandler := localdisk.NewLocalDiskHandler(m.apiClient, recorder)
	ldList, err := ldHandler.ListLocalDisk()
	if err != nil {
		return "", err
	}

	var diskName string
	for _, ld := range ldList.Items {
		if ld.Spec.NodeName == nodeName {
			if ld.Spec.UUID == diskUUID {
				diskName = ld.Spec.DevicePath
				break
			}
		}
	}

	return diskName, nil
}

func (m *manager) updateLocalDisk(replaceDisk *apisv1alpha1.ReplaceDisk) error {
	m.logger.Debug("updateLocalDisk start ... ")
	// replacedDiskName e.g.(/dev/sdb -> sdb)
	diskUuid := replaceDisk.Spec.NewUUID
	nodeName := replaceDisk.Spec.NodeName
	diskName, err := m.getDiskNameByDiskUUID(diskUuid, nodeName)
	if err != nil {
		m.logger.WithError(err).Error("updateLocalDisk: Failed to getDiskNameByDiskUUID")
		return err
	}

	localDisk, err := m.getLocalDiskByDiskName(diskName, nodeName)
	if err != nil {
		m.logger.WithError(err).Error("updateLocalDisk: Failed to getLocalDiskByDiskName")
		return err
	}

	localDisk.Spec.DiskAttributes.Product = ReplaceDiskKind + strconv.Itoa(rand.Intn(1000))
	err = m.localDiskController.UpdateLocalDisk(localDisk)
	if err != nil {
		m.logger.WithError(err).Error("updateLocalDisk: Failed to UpdateLocalDisk")
		return err
	}

	return nil
}
