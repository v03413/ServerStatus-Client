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

func (c *Client) getUpTime(ret *update) {
	uptime, err := host.Uptime()
	if err != nil {
		if c.Debug {

			log.Println(err.Error())
		}

		c.waitGroup.Done()
		return
	}

	ret.Uptime = uptime
	c.waitGroup.Done()
}
func (c *Client) getMemory(ret *update) {
	virtual, err := mem.VirtualMemory()
	if err != nil {
		if c.Debug {

			log.Println(err.Error())
		}

		c.waitGroup.Done()
		return
	}

	ret.MemoryTotal = virtual.Total / 1024
	ret.MemoryUsed = virtual.Used / 1024
	c.waitGroup.Done()
}
func (c *Client) getSwap(ret *update) {
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		if c.Debug {

			log.Println(err.Error())
		}

		c.waitGroup.Done()
		return
	}

	ret.SwapTotal = swapMemory.Total / 1024
	ret.SwapUsed = swapMemory.Used / 1024
	c.waitGroup.Done()
}
func (c *Client) getCpuPercent(ret *update) {
	percent, _ := cpu.Percent(time.Second*time.Duration(c.Interval), false)

	ret.Cpu = percent[0]
	c.waitGroup.Done()
}
func (c *Client) getTraffic(ret *update) {
	var data = traffic{
		in:  0,
		out: 0,
	}
	items, err := psnet.IOCounters(true)
	if err != nil {
		if c.Debug {

			log.Println(err.Error())
		}

		c.waitGroup.Done()
		return
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
		"ip_",
	}

	var contains = func(s string, subStr []string) bool {
		for _, v := range subStr {
			if strings.Contains(s, v) {

				return true
			}
		}

		return false
	}

	for _, info := range items {
		if contains(info.Name, inters) {

			continue
		}

		data.in += info.BytesRecv
		data.out += info.BytesSent
	}

	ret.NetWorkIn = data.in
	ret.NetWorkOut = data.out
	c.waitGroup.Done()
}
