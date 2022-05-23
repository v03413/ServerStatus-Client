package client

import (
	"github.com/shirou/gopsutil/net"
	"log"
	"time"
)

var lastRecvBytes uint64
var lastSendBytes uint64
var lastNetUpdateTime uint64
var currentRx, currentTx uint64

func (c *Client) GetNetRate(retInfo *update) {
	retInfo.NetWorkRx = currentRx
	retInfo.NetWorkTx = currentTx

	c.waitGroup.Done()
}

func (c *Client) startNet() {
	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		var recvTotal, sentTotal uint64
		info, err := net.IOCounters(true)
		if err != nil {
			if c.Debug {

				log.Println(err.Error())
			}

			continue
		}

		for _, v := range info {
			recvTotal += v.BytesRecv
			sentTotal += v.BytesSent
		}

		second := uint64(time.Now().Unix()) - lastNetUpdateTime
		if second > 0 {
			currentRx = (recvTotal - lastRecvBytes) / second
			currentTx = (sentTotal - lastSendBytes) / second
		}

		lastRecvBytes = recvTotal
		lastSendBytes = sentTotal
		lastNetUpdateTime = uint64(time.Now().Unix())
	}
}
