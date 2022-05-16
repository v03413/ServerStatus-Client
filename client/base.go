package client

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net"
	"log"
	"strings"
	"time"
)

func (c *Client) getUpTime() uint64 {
	uptime, err := host.Uptime()
	if err != nil {
		log.Println(err.Error())

		return 0
	}

	return uptime
}
func (c *Client) getMemory() memory {
	virtual, err := mem.VirtualMemory()
	if err != nil {

		return memory{}
	}

	return memory{
		total: virtual.Total / 1024,
		used:  virtual.Used / 1024,
	}
}
func (c *Client) getSwap() swap {
	swapMemory, err := mem.SwapMemory()
	if err != nil {

		return swap{}
	}

	return swap{
		total: swapMemory.Total / 1024,
		used:  swapMemory.Used / 1024,
	}
}
func (c *Client) getCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second*time.Duration(c.Interval), false)

	return percent[0]
}
func (c *Client) getTraffic() traffic {
	var data = traffic{
		in:  0,
		out: 0,
	}
	items, err := psnet.IOCounters(true)
	if err != nil {

		return data
	}

	var inters = []string{
		"lo",
		"tun",
		"docker",
		"veth",
		"br-",
		"vmbr",
		"vnet",
		"kube",
	}

	for _, info := range items {
		for _, v := range inters {
			if !strings.Contains(info.Name, v) {
				continue
			}

			data.in += info.BytesRecv
			data.out += info.BytesSent
		}
	}

	return data
}
