package manager

import (
	"github.com/hwameistor/local-disk-manager/pkg/utils"
	"os"
)

// DiskInfo
type DiskInfo struct {
	// DiskIdentify
	DiskIdentify

	// Attribute
	Attribute Attribute `json:"attribute"`

	// Partition
	Partitions []PartitionInfo `json:"partition"`

	// Raid
	Raid RaidInfo `json:"raid"`
}

// GenerateUUID
func (disk DiskInfo) GenerateUUID() string {
	elementSet := disk.Attribute.Serial + disk.Attribute.Model + disk.Attribute.Vendor + disk.Attribute.WWN

	vitualDiskModels := []string{"EphemeralDisk", "QEMU_HARDDISK", "Virtual_disk"}
	for _, virtualModel := range vitualDiskModels {
		if virtualModel == disk.Attribute.Model {
			host, _ := os.Hostname()
			elementSet += host + disk.Attribute.DevName
		}
	}
	return utils.Hash(elementSet)
}
