package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"strings"
	"sync"

	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
	runtimecache "sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type resources struct {
	apiClient client.Client

	// resourceID is for HA volumes only. Each HA volume must have a unique resourceID.
	// For DRBD, resourceID means the network port.
	// For all non-HA volumes, resourceID is set to '-1'
	allocatedResourceIDs map[string]int
	freeResourceIDList   []int
	maxHAVolumeCount     int

	allocatedStorages *storageCollection
	totalStorages     *storageCollection

	storageNodes map[string]*apisv1alpha1.LocalStorageNode

	lock sync.Mutex

	logger *log.Entry
}

func newResources(maxHAVolumeCount int) *resources {
	return &resources{
		logger:               log.WithField("Module", "Scheduler/Resources"),
		allocatedResourceIDs: make(map[string]int),
		freeResourceIDList:   make([]int, 0, maxHAVolumeCount),
		maxHAVolumeCount:     maxHAVolumeCount,
		allocatedStorages:    newStorageCollection(),
		totalStorages:        newStorageCollection(),
		storageNodes:         map[string]*apisv1alpha1.LocalStorageNode{},
	}
}

func (r *resources) init(apiClient client.Client, informerCache runtimecache.Cache) {
	r.apiClient = apiClient

	// initialize the resources, e.g. resource IDs
	r.initilizeResources()

	nodeInformer, err := informerCache.GetInformer(context.TODO(), &apisv1alpha1.LocalStorageNode{})
	if err != nil {
		r.logger.WithError(err).Fatal("Failed to initiate informer for LocalStorageNode")
	}
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.handleNodeAdd,
		UpdateFunc: r.handleNodeUpdate,
		DeleteFunc: r.handleNodeDelete,
	})

	volumeInformer, err := informerCache.GetInformer(context.TODO(), &apisv1alpha1.LocalVolume{})
	if err != nil {
		r.logger.WithError(err).Fatal("Failed to initiate informer for LocalVolume")
	}
	volumeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: r.handleVolumeUpdate,
	})
}

func (r *resources) initilizeResources() {
	r.logger.Debug("Initializing resources ...")
	volList := &apisv1alpha1.LocalVolumeList{}
	if err := r.apiClient.List(context.TODO(), volList); err != nil {
		r.logger.WithError(err).Fatal("Failed to list LocalVolumes")
	}
	nodeList := &apisv1alpha1.LocalStorageNodeList{}
	if err := r.apiClient.List(context.TODO(), nodeList); err != nil {
		r.logger.WithError(err).Fatal("Failed to list LocalStorageNodes")
	}

	// initialize resource IDs
	usedResourceIDMap := make(map[int]bool)
	for _, vol := range volList.Items {
		if vol.Spec.Config == nil || vol.Spec.Config.ResourceID == -1 || vol.Status.State == apisv1alpha1.VolumeStateDeleted {
			continue
		}
		if !vol.Spec.Config.Convertible && len(vol.Spec.Config.Replicas) < 2 {
			continue
		}
		r.allocatedResourceIDs[vol.Name] = vol.Spec.Config.ResourceID
		usedResourceIDMap[vol.Spec.Config.ResourceID] = true
	}
	for i := 0; i < r.maxHAVolumeCount; i++ {
		if !usedResourceIDMap[i] {
			r.freeResourceIDList = append(r.freeResourceIDList, i)
		}
	}

	// initialize total capacity
	for i := range nodeList.Items {
		r.addTotalStorage(&nodeList.Items[i])
	}
	// initialize allocated capacity
	for i := range volList.Items {
		r.addAllocatedStorage(&volList.Items[i])
	}
}

func (r *resources) predicate(vol *apisv1alpha1.LocalVolume, nodeName string) error {
	node, ok := r.storageNodes[nodeName]
	if !ok {
		return fmt.Errorf("storage node %s not exists", nodeName)
	}
	logCtx := r.logger.WithField("volume", vol.Name)

	totalPool := r.totalStorages.pools[vol.Spec.PoolName]
	allocatedPool := r.allocatedStorages.pools[vol.Spec.PoolName]
	//logCtx.Debug("predicate totalPool = %v, allocatedPool = %v", totalPool, allocatedPool)
	fmt.Printf("predicate totalPool = %+v, allocatedPool = %+v", totalPool, allocatedPool)

	if strings.Contains(strings.Join(vol.Spec.Accessibility.Nodes, ","), nodeName) {
		if vol.Spec.RequiredCapacityBytes > totalPool.capacities[nodeName]-allocatedPool.capacities[nodeName] {
			logCtx.Error("No enough storage capacity on accessibility node")
			return fmt.Errorf("no enough storage capacity on accessibility node")
		}
		if totalPool.volumeCount[nodeName] <= allocatedPool.volumeCount[nodeName] {
			logCtx.Error("No enough volume count on accessibility node")
			return fmt.Errorf("no enough volume count on accessibility node")
		}
		return nil
	}

	if vol.Spec.RequiredCapacityBytes > totalPool.capacities[nodeName]-allocatedPool.capacities[nodeName] {
		return fmt.Errorf("not enough capacity")
	}
	// for disk volume
	if vol.Spec.RequiredCapacityBytes > node.Status.Pools[vol.Spec.PoolName].VolumeCapacityBytesLimit {
		return fmt.Errorf("exceed volume capacity bytes limit")
	}
	if totalPool.volumeCount[nodeName] <= allocatedPool.volumeCount[nodeName] {
		return fmt.Errorf("not enough free volume count")
	}

	return nil
}

// Score calculate node socre for this volume
func (r *resources) Score(vol *apisv1alpha1.LocalVolume, nodeName string) (score int64, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.score(vol, nodeName)
}

func (r *resources) score(vol *apisv1alpha1.LocalVolume, nodeName string) (score int64, err error) {
	if _, ok := r.storageNodes[nodeName]; !ok {
		return 0, fmt.Errorf("storage node %s not exists", nodeName)
	}

	totalPool := r.totalStorages.pools[vol.Spec.PoolName]
	allocatedPool := r.allocatedStorages.pools[vol.Spec.PoolName]
	score = int64(1-float64(vol.Spec.RequiredCapacityBytes)/float64(totalPool.capacities[nodeName]-allocatedPool.capacities[nodeName])) * 100

	return score, nil
}

func (r *resources) getNodeCandidates(vol *apisv1alpha1.LocalVolume) ([]*apisv1alpha1.LocalStorageNode, error) {
	logCtx := r.logger.WithFields(log.Fields{"volume": vol.Name, "spec": vol.Spec})
	logCtx.Debug("getting available nodes for LocalVolume")

	r.lock.Lock()
	defer r.lock.Unlock()

	// step 1. filter out the nodes which have already been allocated
	excludedNodes := map[string]bool{}
	if vol.Spec.Config != nil {
		for _, rep := range vol.Spec.Config.Replicas {
			excludedNodes[rep.Hostname] = true
		}
	}

	// step 2. check the required nodes firstly, if not satisfied, return error immediately
	candidates := []*apisv1alpha1.LocalStorageNode{}
	for _, nn := range vol.Spec.Accessibility.Nodes {
		if len(nn) > 0 && !excludedNodes[nn] {
			if err := r.predicate(vol, nn); err != nil {
				return nil, err
			}
			candidates = append(candidates, r.storageNodes[nn])
			excludedNodes[nn] = true
		}
	}

	// step 3. check the rest of all nodes for the volume replica, and queue the qualified by the avaliable storage space
	pq := make(PriorityQueue, 0)
	for _, node := range r.storageNodes {
		if excludedNodes[node.Name] {
			continue
		}

		if err := r.predicate(vol, node.Name); err != nil {
			continue
		}
		priority, err := r.score(vol, node.Name)
		if err != nil {
			logCtx.Error(err)
			continue
		}
		heap.Push(
			&pq,
			&PriorityItem{
				name:     node.Name,
				priority: priority,
				index:    pq.Len(),
			},
		)
	}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PriorityItem)
		candidates = append(candidates, r.storageNodes[item.name])
		r.logger.WithFields(log.Fields{"node": item.name, "total": pq.Len()}).Debug("Adding a candidate")
	}

	return candidates, nil
}

func (r *resources) getResourceIDForVolume(vol *apisv1alpha1.LocalVolume) (int, error) {
	if vol.Spec.ReplicaNumber < 2 && !vol.Spec.Convertible {
		// try to recycle the resource ID in case of this volume is HA before
		r.recycleResourceID(vol)
		// for non-HA volume, resource ID is -1
		return -1, nil
	}

	return r.allocateResourceID(vol.Name)
}

func (r *resources) allocateResourceID(volName string) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// check if the volume already got resource ID allocated before
	if resID, exists := r.allocatedResourceIDs[volName]; exists {
		return resID, nil
	}

	if len(r.allocatedResourceIDs) >= r.maxHAVolumeCount {
		return -1, fmt.Errorf("can't allocate reourceID, exceeds max volume count")
	}

	if len(r.freeResourceIDList) > 0 {
		resID := r.freeResourceIDList[0]
		r.freeResourceIDList = r.freeResourceIDList[1:]
		r.allocatedResourceIDs[volName] = resID
		return resID, nil
	}

	return -1, fmt.Errorf("can't allocate resource ID")
}

func (r *resources) recycleResourceID(vol *apisv1alpha1.LocalVolume) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if id, exists := r.allocatedResourceIDs[vol.Name]; exists {
		delete(r.allocatedResourceIDs, vol.Name)
		r.freeResourceIDList = append(r.freeResourceIDList, id)
	}
}

func (r *resources) addAllocatedStorage(vol *apisv1alpha1.LocalVolume) {
	if vol.Spec.Config == nil || len(vol.Spec.Config.Replicas) == 0 {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	//r.logger.Debug("addAllocatedStorage = vol.Spec.Config.Replicas = %v", vol.Spec.Config.Replicas)
	fmt.Printf("addAllocatedStorage = vol.Spec.Config.Replicas = %+v", vol.Spec.Config.Replicas)

	for _, replica := range vol.Spec.Config.Replicas {
		// for capacity
		if _, exists := r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname]; !exists {
			r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname] = 0
		}
		r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname] += vol.Spec.Config.RequiredCapacityBytes

		// for volume count
		if _, exists := r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname]; !exists {
			r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname] = 0
		}
		r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname]++
	}
}

func (r *resources) recycleAllocatedStorage(vol *apisv1alpha1.LocalVolume) {
	if vol.Spec.Config == nil || len(vol.Spec.Config.Replicas) == 0 {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	for _, replica := range vol.Spec.Config.Replicas {
		// for capacity
		if _, exists := r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname]; !exists {
			r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname] = 0
		}
		r.allocatedStorages.pools[vol.Spec.PoolName].capacities[replica.Hostname] -= vol.Spec.Config.RequiredCapacityBytes

		// for volume count
		if _, exists := r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname]; !exists {
			r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname] = 0
		}
		r.allocatedStorages.pools[vol.Spec.PoolName].volumeCount[replica.Hostname]--
	}

}

func (r *resources) addTotalStorage(node *apisv1alpha1.LocalStorageNode) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, pool := range node.Status.Pools {
		r.totalStorages.pools[pool.Name].capacities[node.Name] = pool.TotalCapacityBytes
		r.totalStorages.pools[pool.Name].volumeCount[node.Name] = pool.TotalVolumeCount
	}
	r.storageNodes[node.Name] = node
}

func (r *resources) delTotalStorage(node *apisv1alpha1.LocalStorageNode) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, pool := range node.Status.Pools {
		delete(r.totalStorages.pools[pool.Name].capacities, node.Name)
		delete(r.totalStorages.pools[pool.Name].volumeCount, node.Name)
	}
	delete(r.storageNodes, node.Name)
}

func (r *resources) handleNodeAdd(obj interface{}) {
	node := obj.(*apisv1alpha1.LocalStorageNode)
	r.addTotalStorage(node)
}

func (r *resources) handleNodeUpdate(oldObj, newObj interface{}) {
	node := newObj.(*apisv1alpha1.LocalStorageNode)
	r.addTotalStorage(node)
}

func (r *resources) handleNodeDelete(obj interface{}) {
	node := obj.(*apisv1alpha1.LocalStorageNode)
	r.delTotalStorage(node)

}

func (r *resources) handleVolumeUpdate(oldObj, newObj interface{}) {
	oVol := oldObj.(*apisv1alpha1.LocalVolume)
	nVol := newObj.(*apisv1alpha1.LocalVolume)

	// 1. calculate allocated capacity according to LocalVolume.Spec.Config
	// recycle old volume
	r.recycleAllocatedStorage(oVol)
	// allocate new volume
	r.addAllocatedStorage(nVol)

	// 2. recycle resource ID when LocalVolume is deleted
	if nVol.Status.State == apisv1alpha1.VolumeStateDeleted {
		r.recycleResourceID(nVol)
	}
	if nVol.Spec.Config == nil {
		r.recycleResourceID(nVol)
	} else if !nVol.Spec.Config.Convertible && len(nVol.Spec.Config.Replicas) < 2 {
		r.recycleResourceID(nVol)
	}

}
