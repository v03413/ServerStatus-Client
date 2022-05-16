package client

import (
	"github.com/shirou/gopsutil/disk"
	"sync"
	"time"
)

var validFs = []string{"ext4", "ext3", "ext2", "reiserfs", "overlay", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "apfs", "hfs"}
var lastReadBytes uint64
var lastWriteBytes uint64
var ioLock sync.Mutex

func (c *Client) getDiskIo() ioStat {
	var read, write uint64

	stats, err := disk.IOCounters()
	if err != nil {

		return ioStat{}
	}

	for _, stat := range stats {
		read += stat.ReadBytes
		write += stat.WriteBytes
	}

	var ret = ioStat{
		writeBytes: write - lastWriteBytes,
		readBytes:  read - lastReadBytes,
		second:     uint64(time.Now().Sub(c.lastUpdateTime).Seconds()),
	}

	ioLock.Lock()
	lastReadBytes = read
	lastWriteBytes = write
	ioLock.Unlock()

	return ret
}

func (c *Client) getDiskUsage() diskStat {
	var disks []string
	var size uint64
	var used uint64

	list, err := disk.Partitions(true)
	if err != nil {

		return diskStat{}
	}

	for _, itm := range list {
		for _, valid := range validFs {
			if itm.Fstype == valid {

				disks = append(disks, itm.Mountpoint)
			}
		}
	}

	for _, point := range disks {
		if stat, err := disk.Usage(point); err == nil {
			size += stat.Total
			used += stat.Used
		}
	}

	return diskStat{
		size: size / 1024 / 1024,
		used: used / 1024 / 1024,
	}
}
