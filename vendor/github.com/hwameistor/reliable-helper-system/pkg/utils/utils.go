package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	lsapisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
)

// disk class
const (
	DiskClassNameHDD  = "HDD"
	DiskClassNameSSD  = "SSD"
	DiskClassNameNVMe = "NVMe"
)

// consts
const (
	PoolNamePrefix  = "LocalStorage_Pool"
	PoolNameForHDD  = PoolNamePrefix + DiskClassNameHDD
	PoolNameForSSD  = PoolNamePrefix + DiskClassNameSSD
	PoolNameForNVMe = PoolNamePrefix + DiskClassNameNVMe
)

// GetNodeName gets the node name from env, else
// returns an error
func GetNodeName() string {
	nodeName, ok := os.LookupEnv("NODENAME")
	if !ok {
		log.Errorf("Failed to get NODENAME from ENV")
		return ""
	}

	return nodeName
}

// GetNamespace get Namespace from env, else it returns error
func GetNamespace() string {
	ns, ok := os.LookupEnv("NAMESPACE")
	if !ok {
		log.Errorf("Failed to get NameSpace from ENV")
		return ""
	}

	return ns
}

func MergeLocalVolumeReplicaMap(localVolumeReplicaMap ...map[string][]*lsapisv1alpha1.LocalVolumeReplica) map[string][]*lsapisv1alpha1.LocalVolumeReplica {
	newLocalVolumeReplicaMap := map[string][]*lsapisv1alpha1.LocalVolumeReplica{}
	for _, m := range localVolumeReplicaMap {
		for k, v := range m {
			newLocalVolumeReplicaMap[k] = v
		}
	}
	return newLocalVolumeReplicaMap
}

// ConvertNodeName e.g.(10.23.10.12 => 10-23-10-12)
func ConvertNodeName(node string) string {
	return strings.Replace(node, ".", "-", -1)
}

func GetPoolNameAccordingDiskType(diskType string) (string, error) {
	switch diskType {
	case DiskClassNameHDD:
		return PoolNameForHDD, nil
	case DiskClassNameSSD:
		return PoolNameForSSD, nil
	case DiskClassNameNVMe:
		return PoolNameForNVMe, nil
	}
	return "", fmt.Errorf("not supported pool type %s", diskType)
}
