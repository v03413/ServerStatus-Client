package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/load"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const Version = 0.15
const DefaultServer = "127.0.0.1"
const DefaultPort = 35601
const DefaultInterval = 1
const DefaultUsername = "s01"
const DefaultPassword = "USER_DEFAULT_PASSWORD"
const DefaultProtocol = "ip4"
const PingPacketHistoryLen = 100
const TimeOut = time.Second * 5

const ProbePort = 80
const PingCu = "cu.tz.cloudcpp.com"
const PingCt = "ct.tz.cloudcpp.com"
const PingCm = "cm.tz.cloudcpp.com"

type Client struct {
	Server    string
	Port      uint64
	Username  string
	Password  string
	Interval  uint64
	Protocol  string
	Debug     bool
	waitGroup sync.WaitGroup
	conn      net.Conn
	lastTime  time.Time
	pingTime  sync.Map
	baseInfo  struct {
		checkIp uint8
		timer   uint8
	}
}

func (c *Client) Start() error {
	if err := c.initiation(); err != nil {

		return errors.New(fmt.Sprintf("初始化错误：%s", err.Error()))
	}
	if err := c.connectServer(); err != nil {

		return errors.New(err.Error())
	}

	go c.startPing()
	go c.startNet()
	go c.startDiskIo()
	go c.startRun()

	return nil
}
func (c *Client) startRun() {
	defer func(conn net.Conn) {
		_ = conn.Close()

	}(c.conn)

	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		var start = time.Now()
		var update = c.getUpdateInfo()
		if c.Debug {
			log.Printf("获取耗时：%v", time.Now().Sub(start))
		}
		var data = []byte("update ")
		if jsonByte, err := json.Marshal(update); err != nil {
			log.Println(err.Error())
		} else {
			_ = c.conn.SetWriteDeadline(time.Now().Add(TimeOut))
			data = append(data, jsonByte...)
			data = append(data, []byte("\n")...)
			write, err := c.conn.Write(data)
			if err != nil {
				_ = c.conn.Close()
				log.Printf("[准备重连]发送失败：%s\n", err.Error())

				if err = c.connectServer(); err != nil {

					log.Printf("服务器重连失败：%s\n", err.Error())
				}
			} else {
				log.Printf("发送成功：%dByte\n", write)
			}
		}
	}
}
func (c *Client) connectServer() error {
	var recvData = make([]byte, 128)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%v", c.Server, c.Port), TimeOut)
	if err != nil {

		return errors.New(fmt.Sprintf("[连接]建立失败：%s", err.Error()))
	}

	for {
		_ = conn.SetReadDeadline(time.Now().Add(TimeOut))
		_ = conn.SetWriteDeadline(time.Now().Add(TimeOut))
		_, err = conn.Read(recvData)
		if err != nil {

			return errors.New(fmt.Sprintf("[连接]数据响应错误：%s，%s", err.Error(), string(recvData)))
		}

		if strings.Contains(string(recvData), "Authentication required") {
			_, err := conn.Write([]byte(fmt.Sprintf("%s:%s\n", c.Username, c.Password)))
			if err != nil {

				return errors.New(fmt.Sprintf("[连接]数据发送失败：%s", err.Error()))
			}

			continue
		}

		if strings.Contains(string(recvData), "You are connecting via") {
			if !strings.Contains(string(recvData), "IPv4") {
				c.baseInfo.checkIp = 4
			} else {
				c.baseInfo.checkIp = 6
			}

			break
		}
	}

	c.conn = conn

	log.Println("服务器连接成功")

	return nil
}
func (c *Client) initiation() error {
	if c.Server == "" {

		c.Server = DefaultServer
	}
	if c.Username == "" {

		c.Username = DefaultUsername
	}
	if c.Password == "" {

		c.Password = DefaultPassword
	}
	if c.Protocol == "" {

		c.Protocol = DefaultProtocol
	}
	if c.Port == 0 {

		c.Port = DefaultPort
	}
	if c.Interval <= 0 {

		c.Interval = DefaultInterval
	}

	c.pingTime = sync.Map{}

	return nil
}
func (c *Client) getUpdateInfo() update {
	c.waitGroup = sync.WaitGroup{}
	var ret = &update{}

	c.waitGroup.Add(1)
	go c.getUpTime(ret)

	c.waitGroup.Add(1)
	go c.getCpuPercent(ret)

	c.waitGroup.Add(1)
	go c.getMemory(ret)

	c.waitGroup.Add(1)
	go c.getSwap(ret)

	c.waitGroup.Add(1)
	go c.getDiskUsage(ret)

	c.waitGroup.Add(1)
	go c.getTraffic(ret)

	c.waitGroup.Add(1)
	go c.GetNetRate(ret)

	ret.Ping10086 = c.getLostPacket("10086")
	ret.Ping10010 = c.getLostPacket("10010")
	ret.Ping189 = c.getLostPacket("189")

	ret.Time10086 = c.getPingTime("10086")
	ret.Time10010 = c.getPingTime("10010")
	ret.Time189 = c.getPingTime("189")

	c.waitGroup.Add(1)
	go c.getDiskIo(ret)

	c.waitGroup.Add(1)
	go c.getTupd(ret)

	if loadavg, err := load.Avg(); err == nil {
		ret.Load1 = loadavg.Load1
		ret.Load5 = loadavg.Load5
		ret.Load15 = loadavg.Load15
	}

	ret.IpStatus = true
	c.waitGroup.Wait()

	return *ret
}
