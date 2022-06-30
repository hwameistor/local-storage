package node

import (
	"context"
	"fmt"
	"github.com/hwameistor/local-storage/pkg/member/node/storage"
	"k8s.io/client-go/tools/record"
	"net"
	"os"

	ldmv1alpha1 "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	localdisk2 "github.com/hwameistor/local-disk-manager/pkg/handler/localdisk"
	"github.com/hwameistor/local-storage/pkg/apis"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-storage/pkg/common"
	"github.com/hwameistor/local-storage/pkg/member/node/diskmonitor"
	rdmgr "github.com/hwameistor/reliable-helper-system/pkg/replacedisk/manager"
	log "github.com/sirupsen/logrus"

	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	runtimecache "sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// maxRetries is the number of times a task will be retried before it is dropped out of the queue.
// With the current rate-limiter in use math.Max(16s, (1s*2^(maxRetries-1))) the following numbers represent the times
// a task is going to be requeued:
//
// Infinitely retry
const maxRetries = 0

type manager struct {
	name string

	namespace string

	apiClient client.Client

	informersCache runtimecache.Cache

	// to record all the replicas located at this node, volumeName -> replicaName
	replicaRecords map[string]string

	storageMgr *storage.LocalManager

	// if there is any suspicious volume replica, put it in this queue to check health
	// for example, when a disk runs into problem, the associated volume replicas should be added into this queue
	//	healthCheckQueue *common.TaskQueue

	diskEventQueue *diskmonitor.EventQueue

	volumeTaskQueue *common.TaskQueue

	volumeReplicaTaskQueue *common.TaskQueue

	localDiskClaimTaskQueue *common.TaskQueue

	localDiskTaskQueue *common.TaskQueue

	configManager *configManager

	ldhandler *localdisk2.LocalDiskHandler

	rdhandler *rdmgr.ReplaceDiskHandler

	logger *log.Entry
}

// New node manager
func New(name string, namespace string, cli client.Client, informersCache runtimecache.Cache, config apisv1alpha1.SystemConfig) (apis.NodeManager, error) {
	var recorder record.EventRecorder
	configManager, err := NewConfigManager(name, config, cli)
	if err != nil {
		return nil, err
	}
	return &manager{
		name:                    name,
		namespace:               namespace,
		apiClient:               cli,
		informersCache:          informersCache,
		replicaRecords:          map[string]string{},
		volumeTaskQueue:         common.NewTaskQueue("VolumeTask", maxRetries),
		volumeReplicaTaskQueue:  common.NewTaskQueue("VolumeReplicaTask", maxRetries),
		localDiskClaimTaskQueue: common.NewTaskQueue("LocalDiskClaim", maxRetries),
		localDiskTaskQueue:      common.NewTaskQueue("LocalDisk", maxRetries),
		// healthCheckQueue:        common.NewTaskQueue("HealthCheckTask", maxRetries),
		diskEventQueue: diskmonitor.NewEventQueue("DiskEvents"),
		configManager:  configManager,
		ldhandler:      localdisk2.NewLocalDiskHandler(cli, recorder),
		rdhandler:      rdmgr.NewReplaceDiskHandler(cli, recorder),
		logger:         log.WithField("Module", "NodeManager"),
	}, nil
}

func (m *manager) Run(stopCh <-chan struct{}) {
	m.initCache()

	m.register()

	m.setupInformers()

	go m.startVolumeTaskWorker(stopCh)

	go m.startVolumeReplicaTaskWorker(stopCh)

	go m.startLocalDiskClaimTaskWorker(stopCh)

	go m.startLocalDiskTaskWorker(stopCh)

	go m.startDiskEventWorker(stopCh)

	go diskmonitor.New(m.diskEventQueue).Run(stopCh)

	go m.configManager.Run(stopCh)

	// move disk health check out, as a separate process
	//go healths.NewDiskHealthManager(m.name, m.apiClient).Run(stopCh)
}

/*
func (m *manager) isPhysicalNode() bool {
	params := exechelper.ExecParams{
		CmdName: "cat",
		CmdArgs: []string{"/sys/class/dmi/id/product_name"},
	}
	res := basicexecutor.New().RunCommand(params)
	if res.ExitCode != 0 {
		m.logger.WithError(res.Error).Fatal("Can't determine if the node is physical or virtual")
	}
	virtualVendors := []string{
		"vmware",
		"kvm",
		"virtualbox",
		"qemu",
	}
	vendorStr := strings.ToLower(res.OutBuf.String())
	for _, vv := range virtualVendors {
		if strings.Contains(vendorStr, vv) {
			return false
		}
	}
	return true
}
*/

func (m *manager) initCache() {
	// initialize replica records
	m.logger.Debug("Initializing replica records in cache")
	replicaList := &apisv1alpha1.LocalVolumeReplicaList{}
	if err := m.apiClient.List(context.TODO(), replicaList); err != nil {
		m.logger.WithError(err).Fatal("Failed to list replicas")
	}
	for _, replica := range replicaList.Items {
		if replica.Spec.NodeName == m.name {
			m.replicaRecords[replica.Spec.VolumeName] = replica.Name
		}
	}
}

func (m *manager) setupInformers() {
	nodeInformer, err := m.informersCache.GetInformer(context.TODO(), &apisv1alpha1.LocalStorageNode{})
	if err != nil {
		// error happens, crash the node
		m.logger.WithError(err).Fatal("Failed to get informer for Node")
	}
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// protect from being deleted by mistake
		DeleteFunc: m.handleNodeDelete,
	})

	volumeReplicaInformer, err := m.informersCache.GetInformer(context.TODO(), &apisv1alpha1.LocalVolumeReplica{})
	if err != nil {
		// error happens, crash the node
		m.logger.WithError(err).Fatal("Failed to get informer for VolumeReplica")
	}
	volumeReplicaInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// protect from being deleted by mistake
		DeleteFunc: m.handleVolumeReplicaDelete,
		// for updating local storage node status for volume replica
		UpdateFunc: m.handleVolumeReplicaUpdate,
	})

	localDiskClaimInformer, err := m.informersCache.GetInformer(context.TODO(), &ldmv1alpha1.LocalDiskClaim{})
	if err != nil {
		// error happens, crash the node
		//m.logger.WithError(err).Fatal("Failed to get informer for LocalDiskClaim")
		m.logger.WithError(err).Fatal("Failed to get informer for LocalDiskClaim")
	}
	localDiskClaimInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    m.handleLocalDiskClaimAdd,
		UpdateFunc: m.handleLocalDiskClaimUpdate,
	})

	localDiskInformer, err := m.informersCache.GetInformer(context.TODO(), &ldmv1alpha1.LocalDisk{})
	if err != nil {
		// error happens, crash the node
		m.logger.WithError(err).Fatal("Failed to get informer for LocalDisk")
	}
	localDiskInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    m.handleLocalDiskAdd,
		UpdateFunc: m.handleLocalDiskUpdate,
	})

}

func (m *manager) Storage() *storage.LocalManager {
	return m.storageMgr
}

func (m *manager) TakeVolumeReplicaTaskAssignment(vol *apisv1alpha1.LocalVolume) {
	// have to add all volumes into the assignment queue, even this node is not in volume.config
	// in case of removing replica, it is not in the volume.config but should be recycled
	m.volumeTaskQueue.Add(vol.Name)
}

func (m *manager) ReconcileVolumeReplica(replica *apisv1alpha1.LocalVolumeReplica) {
	if replica.Spec.NodeName == m.name {
		m.volumeReplicaTaskQueue.Add(replica.Name)
	}
}

func (m *manager) register() {
	var nodeConfig *apisv1alpha1.NodeConfig
	logCtx := m.logger.WithFields(log.Fields{"node": m.name})
	logCtx.Debug("Registering node into cluster")
	k8sNode := &k8scorev1.Node{}
	if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Name: m.name}, k8sNode); err != nil {
		logCtx.WithError(err).Fatal("Can't find K8S node")
	}

	myNode := &apisv1alpha1.LocalStorageNode{}
	if err := m.apiClient.Get(context.TODO(), types.NamespacedName{Name: m.name}, myNode); err != nil {
		if !errors.IsNotFound(err) {
			logCtx.WithError(err).Fatal("Failed to get Node info")
		}
		myNode.Name = m.name
		myNode.Spec.HostName = m.name
		nodeConfig, err = m.getConfByK8SNodeOrDefault(k8sNode)
		if err != nil {
			logCtx.WithError(err).Fatal("Failed to get Node configuration")
		}
		if err := m.configNode(nodeConfig, myNode); err != nil {
			logCtx.WithError(err).Fatal("Failed to config node when register node.")
		}
		if err = m.apiClient.Create(context.TODO(), myNode); err != nil {
			logCtx.WithError(err).Fatal("Can not create Node when registering.")
		}
	} else {
		if len(myNode.Spec.StorageIP) == 0 {
			// for upgrade
			ipAddr, err := m.getStorageIPv4Address(k8sNode)
			if err != nil {
				logCtx.WithError(err).Fatal("Failed to get IPv4 address")
			}
			myNode.Spec.StorageIP = ipAddr
			if err = m.apiClient.Update(context.TODO(), myNode); err != nil {
				logCtx.WithError(err).Fatal("Failed to update Kubernetes Node for IP address")
			}
		}
		nodeConfig = m.getNodeConf(myNode)
	}
	nodeConfig.Name = m.name

	m.storageMgr = storage.NewLocalManager(nodeConfig, m.apiClient)
	if err := m.storageMgr.Register(); err != nil {
		logCtx.WithError(err).Fatal("Failed to register node's storage manager")
	}
}

func (m *manager) getNodeConf(node *apisv1alpha1.LocalStorageNode) *apisv1alpha1.NodeConfig {
	return &apisv1alpha1.NodeConfig{
		StorageIP: node.Spec.StorageIP,
		Topology:  node.Spec.Topo.DeepCopy(),
	}
}

func (m *manager) configNode(config *apisv1alpha1.NodeConfig, node *apisv1alpha1.LocalStorageNode) error {
	if config.Topology != nil {
		node.Spec.Topo = *config.Topology
	}
	node.Spec.StorageIP = config.StorageIP
	return nil
}

func (m *manager) getConfByK8SNodeOrDefault(k8sNode *k8scorev1.Node) (*apisv1alpha1.NodeConfig, error) {
	ipAddr, err := m.getStorageIPv4Address(k8sNode)
	if err != nil {
		return nil, err
	}
	return &apisv1alpha1.NodeConfig{StorageIP: ipAddr}, nil

}

func (m *manager) getStorageIPv4Address(k8sNode *k8scorev1.Node) (string, error) {
	logCtx := m.logger.WithField("node", k8sNode.Name)
	// lookup from k8s node's annotation firstly
	annotationKey := os.Getenv(apisv1alpha1.StorageIPv4AddressAnnotationKeyEnv)
	if len(annotationKey) > 0 {
		ipAddr, has := k8sNode.Annotations[annotationKey]
		if has {
			if net.ParseIP(ipAddr) != nil {
				return ipAddr, nil
			}
			logCtx.WithFields(log.Fields{"annotationKey": annotationKey, "ip": ipAddr}).Error("Invalid IPv4 address")
			return "", fmt.Errorf("invalid IPv4 address")
		}
		logCtx.WithField("annotationKey", annotationKey).Info("Not found in Kubernetes Node")
	}

	// lookup from k8s node's addresses
	for _, addr := range k8sNode.Status.Addresses {
		if addr.Type == k8scorev1.NodeInternalIP {
			return addr.Address, nil
		}
	}

	return "", fmt.Errorf("not found valid IPv4 address")
}

func (m *manager) handleVolumeReplicaUpdate(oldObj, newObj interface{}) {
	replica, _ := newObj.(*apisv1alpha1.LocalVolumeReplica)
	if replica.Spec.NodeName != m.name {
		return
	}
	m.storageMgr.UpdateNodeForVolumeReplica(replica)
}

func (m *manager) handleLocalDiskClaimUpdate(oldObj, newObj interface{}) {
	localDiskClaim, _ := newObj.(*ldmv1alpha1.LocalDiskClaim)
	if localDiskClaim.Spec.NodeName != m.name {
		return
	}
	m.localDiskClaimTaskQueue.Add(localDiskClaim.Namespace + "/" + localDiskClaim.Name)
}

func (m *manager) handleLocalDiskClaimAdd(obj interface{}) {
	localDiskClaim, _ := obj.(*ldmv1alpha1.LocalDiskClaim)
	if localDiskClaim.Spec.NodeName != m.name {
		return
	}
	m.localDiskClaimTaskQueue.Add(localDiskClaim.Namespace + "/" + localDiskClaim.Name)
}

func (m *manager) handleLocalDiskUpdate(oldObj, newObj interface{}) {
	localDisk, _ := newObj.(*ldmv1alpha1.LocalDisk)
	if localDisk.Spec.NodeName != m.name {
		return
	}
	m.localDiskTaskQueue.Add(localDisk.Namespace + "/" + localDisk.Name)
}

func (m *manager) handleLocalDiskAdd(obj interface{}) {
	localDisk, _ := obj.(*ldmv1alpha1.LocalDisk)
	if localDisk.Spec.NodeName != m.name {
		return
	}
	m.localDiskTaskQueue.Add(localDisk.Namespace + "/" + localDisk.Name)
}

func (m *manager) handleVolumeReplicaDelete(obj interface{}) {
	replica, _ := obj.(*apisv1alpha1.LocalVolumeReplica)
	if replica.Spec.NodeName != m.name {
		return
	}

	m.logger.WithFields(log.Fields{"replica": replica.Name}).Info("Observed a VolumeReplica CRD deletion...")
	if replica.Status.State != apisv1alpha1.VolumeReplicaStateDeleted {
		// must be deleted by a mistake, rebuild it
		m.logger.WithFields(log.Fields{"replica": replica.Name, "spec": replica.Spec, "status": replica.Status}).Warning("Rebuilding VolumeReplica CRD ...")
		// TODO: need retry considering the case of creating failure??
		newReplica := &apisv1alpha1.LocalVolumeReplica{}
		newReplica.Name = replica.Name
		newReplica.Spec = replica.Spec

		if err := m.apiClient.Create(context.TODO(), newReplica); err != nil {
			m.logger.WithFields(log.Fields{"replica": replica.Name}).WithError(err).Error("Failed to rebuild VolumeReplica")
		}
	} else {
		delete(m.replicaRecords, replica.Spec.VolumeName)
	}
}

func (m *manager) handleNodeDelete(obj interface{}) {
	node, _ := obj.(*apisv1alpha1.LocalStorageNode)
	if node.Name != m.name {
		return
	}
	m.logger.WithFields(log.Fields{"node": node.Name}).Info("Observed a Node CRD deletion...")

	// must be deleted by a mistake, rebuild it
	m.logger.Warning("Rebuilding Node CRD ...")
	// TODO: need retry considering the case of creating failure??
	nodeToRecovery := &apisv1alpha1.LocalStorageNode{}
	nodeToRecovery.SetName(node.GetName())
	nodeToRecovery.Spec = node.Spec
	nodeToRecovery.Status = node.Status
	if err := m.apiClient.Create(context.TODO(), nodeToRecovery); err != nil {
		m.logger.WithFields(log.Fields{"node": nodeToRecovery.GetName()}).WithError(err).Error("Failed to rebuild VolumeReplica")
	}
	if err := m.apiClient.Status().Update(context.TODO(), nodeToRecovery); err != nil {
		m.logger.WithFields(log.Fields{"node": nodeToRecovery.GetName()}).WithError(err).Error("Failed to rebuild VolumeReplica")
	}
}
