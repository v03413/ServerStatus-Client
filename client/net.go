package client

import (
	"github.com/shirou/gopsutil/net"
	"sync"
	"time"
)

var lastRecvBytes uint64
var lastSendBytes uint64
var netLock sync.Mutex

func (c *Client) GetNetRate() rateStat {
	var recvTotal, sentTotal uint64
	info, _ := net.IOCounters(true)
	for _, v := range info {
		recvTotal += v.BytesRecv
		sentTotal += v.BytesSent
	}

	var ret = rateStat{
		recvBytes: recvTotal - lastRecvBytes,
		sendBytes: sentTotal - lastSendBytes,
		second:    uint64(time.Now().Sub(c.lastUpdateTime).Seconds()),
	}

	netLock.Lock()
	lastRecvBytes = recvTotal
	lastSendBytes = sentTotal
	netLock.Unlock()

	return ret
}
