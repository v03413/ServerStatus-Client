package client

import (
	"github.com/go-ping/ping"
	"log"
	"sync"
	"time"
)

var CmLostPacket = newLostPacket(PingPacketHistoryLen) // 移动
var CuLostPacket = newLostPacket(PingPacketHistoryLen) // 联通
var CtLostPacket = newLostPacket(PingPacketHistoryLen) // 电信
var pingHost = sync.Map{}

func init() {
	pingHost.Store("cu", PingCu)
	pingHost.Store("cm", PingCm)
	pingHost.Store("ct", PingCt)
}

func (c *Client) getPingTime(host string) uint {
	if v, ok := c.pingTime.Load(host); ok {

		return v.(uint)
	}

	return 0
}

func (c *Client) startPing() {
	timeout := time.Second * time.Duration(c.Interval)

	for range time.Tick(timeout) {
		pingHost.Range(func(k, v interface{}) bool {
			var lost *lostPacket

			isp := k.(string)
			domain := v.(string)

			if isp == "cm" {
				lost = &CmLostPacket
			} else if isp == "cu" {
				lost = &CuLostPacket
			} else if isp == "ct" {
				lost = &CtLostPacket
			} else {
				return false
			}

			pinger, err := ping.NewPinger(domain)
			if err != nil {

				log.Println(err.Error())
				return false
			}

			pinger.Timeout = timeout
			pinger.Count = 1
			pinger.SetPrivileged(true)

			if c.Protocol == DefaultProtocol {
				pinger.SetNetwork("ipv4")
			} else {
				pinger.SetNetwork("ipv6")
			}

			err = pinger.Run()
			if err != nil {
				lost.Push(true)

				log.Println(err.Error())
				return false
			}

			stat := pinger.Statistics()
			if stat.AvgRtt == 0 {

				lost.Push(true)
				c.pingTime.Store(isp, uint(time.Duration(0)))
				return true
			}

			lost.Push(false)
			c.pingTime.Store(isp, uint(stat.AvgRtt/time.Millisecond))

			return true
		})
	}
}

func (c *Client) getLostPacket(isp string) float64 {
	if isp == "cm" {

		return CmLostPacket.Get()
	}
	if isp == "cu" {

		return CuLostPacket.Get()
	}
	if isp == "ct" {

		return CtLostPacket.Get()
	}

	return 0
}

type lostPacket struct {
	size uint64
	lock sync.RWMutex
	data []bool
}

func (p *lostPacket) Push(v bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.data) == int(p.size) {

		p.data = p.data[1:]
	}

	p.data = append(p.data, v)
}

func (p *lostPacket) Get() float64 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var succ, total uint64
	for _, v := range p.data {
		total += 1
		if v {

			succ += 1
		}
	}

	if total == 0 {

		return 0
	}

	// 这里故意加 0.001，是由于云端的一些Bug；如果直接返回整数，云端无法读取
	return float64(succ)/float64(total)*100 + 0.001
}

func newLostPacket(size uint64) lostPacket {
	return lostPacket{
		size: size,
	}
}
