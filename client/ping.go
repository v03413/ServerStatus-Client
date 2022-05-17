package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

var lostPacket = make(map[string][]bool)
var pingHost = make(map[string]string)

func init() {
	lostPacket["10010"] = []bool{}
	lostPacket["10086"] = []bool{}
	lostPacket["189"] = []bool{}

	pingHost["10010"] = PingCu
	pingHost["10086"] = PingCm
	pingHost["189"] = PingCt
}

func (c *Client) getPingTime(host string) uint {
	if v, ok := c.pingTime[host]; ok {

		return v
	}

	return 0
}
func (c *Client) startPing() {
	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		for host, domain := range pingHost {
			var ip *net.IPAddr

			if c.Protocol == DefaultProtocol {
				ip, _ = net.ResolveIPAddr(DefaultProtocol, domain)
			} else {
				ip, _ = net.ResolveIPAddr("ip6", domain)
			}

			if len(lostPacket[host]) >= PingPacketHistoryLen {
				// 弹出第一个

				lostPacket[host] = lostPacket[host][1:]
			}

			var start = time.Now()
			dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, ProbePort))
			if err == nil {
				_ = dial.Close()

				lostPacket[host] = append(lostPacket[host], false)
				c.pingTime[host] = uint(time.Now().Sub(start).Milliseconds())
				continue
			}

			lostPacket[host] = append(lostPacket[host], true)
			log.Println(err.Error())
		}
	}
}
func (c *Client) getLostPacket(host string) float64 {
	var succ, total uint64
	for _, v := range lostPacket[host] {
		total += 1
		if v {
			succ += 1
		}
	}

	if total == 0 {

		return 0
	}

	return float64(succ / total)
}
