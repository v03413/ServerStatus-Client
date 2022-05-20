package client

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var lostPacket10086 = newLostPacket(PingPacketHistoryLen)
var lostPacket10010 = newLostPacket(PingPacketHistoryLen)
var lostPacket189 = newLostPacket(PingPacketHistoryLen)
var pingHost = sync.Map{}

func init() {
	pingHost.Store("10010", PingCu)
	pingHost.Store("10086", PingCm)
	pingHost.Store("189", PingCt)
}

func (c *Client) getPingTime(host string) uint {
	if v, ok := c.pingTime[host]; ok {

		return v
	}

	return 0
}

func (c *Client) startPing() {
	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		pingHost.Range(func(k, v interface{}) bool {
			var ip *net.IPAddr
			var lost *lostPacket

			host := k.(string)
			domain := v.(string)

			if c.Protocol == DefaultProtocol {
				ip, _ = net.ResolveIPAddr(DefaultProtocol, domain)
			} else {
				ip, _ = net.ResolveIPAddr("ip6", domain)
			}

			if host == "10086" {
				lost = &lostPacket10086
			} else if host == "10010" {
				lost = &lostPacket10010
			} else {
				lost = &lostPacket189
			}

			var start = time.Now()
			dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, ProbePort))
			if err == nil {
				_ = dial.Close()

				lost.Push(false)
				c.pingTime[host] = uint(time.Now().Sub(start).Milliseconds())
			} else {
				lost.Push(true)
			}

			return true
		})
	}
}

func (c *Client) getLostPacket(host string) float64 {
	if host == "10086" {

		return lostPacket10086.Get()
	}
	if host == "10010" {

		return lostPacket10010.Get()
	}
	if host == "189" {

		return lostPacket189.Get()
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
