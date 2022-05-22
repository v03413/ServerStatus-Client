package client

import (
	"github.com/shirou/gopsutil/disk"
	"log"
	"time"
)

var validFs = []string{"ext4", "ext3", "ext2", "reiserfs", "overlay", "jfs", "btrfs", "fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "apfs", "hfs"}
var lastReadBytes uint64
var lastWriteBytes uint64
var lastDiskUpdateTime uint64
var curReadIo, curWriteIo uint64

func (c *Client) getDiskIo(retInfo *update) {

	retInfo.IoRead = curReadIo
	retInfo.IoWrite = curWriteIo

	c.waitGroup.Done()
}

func (c *Client) getDiskUsage(ret *update) {
	var disks []string
	var size uint64
	var used uint64

	list, err := disk.Partitions(true)
	if err != nil {
		if c.Debug {

			log.Println(err.Error())
		}

		c.waitGroup.Done()
		return
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

	ret.HddTotal = size / 1024 / 1024
	ret.HddUsed = used / 1024 / 1024
	c.waitGroup.Done()
}

func (c *Client) startDiskIo() {
	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		var read, write uint64

		stats, err := disk.IOCounters()
		if err != nil {
			if c.Debug {

				log.Println(err.Error())
			}

			continue
		}

		for _, stat := range stats {
			read += stat.ReadBytes
			write += stat.WriteBytes
		}

		second := uint64(time.Now().Unix()) - lastDiskUpdateTime
		if second > 0 {
			curReadIo = (read - lastReadBytes) / second
			curWriteIo = (write - lastWriteBytes) / second
		}

		lastReadBytes = read
		lastWriteBytes = write
		lastDiskUpdateTime = uint64(time.Now().Unix())
	}
}
