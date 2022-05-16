package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/load"
	"log"
	"net"
	"strings"
	"time"
)

const DefaultServer = "127.0.0.1"
const DefaultPort = 35601
const DefaultInterval = 1
const DefaultUsername = "s01"
const DefaultPassword = "USER_DEFAULT_PASSWORD"
const DefaultProtocol = "ip4"
const PingPacketHistoryLen = 100

const ProbePort = 80
const PingCu = "cu.tz.cloudcpp.com"
const PingCt = "ct.tz.cloudcpp.com"
const PingCm = "cm.tz.cloudcpp.com"

type Client struct {
	Server   string
	Port     uint64
	Username string
	Password string
	Interval uint64
	Protocol string
	conn     net.Conn
	baseInfo struct {
		checkIp uint8
		timer   uint8
	}
	lastUpdateTime time.Time
	pingTime       map[string]uint
}

func (c *Client) Start() error {
	if err := c.initiation(); err != nil {

		return err
	}
	if err := c.connectServer(); err != nil {

		return err
	}

	log.Println("服务器授权成功")

	go c.startPing()
	go c.startRun()

	return nil
}
func (c *Client) startRun() {
	defer func(conn net.Conn) {
		_ = conn.Close()

	}(c.conn)

	for range time.Tick(time.Second * time.Duration(c.Interval)) {
		var update = c.getUpdateInfo()
		var data = []byte("update ")
		if jsonByte, err := json.Marshal(update); err != nil {
			log.Println(err.Error())
		} else {
			data = append(data, jsonByte...)
			data = append(data, []byte("\n")...)
			write, err := c.conn.Write(data)
			if err != nil {
				log.Printf("发送失败：%s\n", err.Error())
			} else {
				log.Printf("发送成功：%dByte\n", write)
			}
		}
	}
}
func (c *Client) connectServer() error {
	var recvData = make([]byte, 1024)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", c.Server, c.Port))
	if err != nil {

		return err
	}

	_, err = conn.Read(recvData)
	if err != nil {

		return err
	}

	if strings.Contains(string(recvData), "Authentication required") {
		_, err := conn.Write([]byte(fmt.Sprintf("%s:%s\n", c.Username, c.Password)))
		if err != nil {

			return err
		}

		_, err = conn.Read(recvData)
		if err != nil {

			return err
		}

		if !strings.HasPrefix(string(recvData), "Authentication successful") {

			return errors.New("服务器拒绝授权")
		}

		_, err = conn.Read(recvData)
		if err != nil {

			return err
		}

		if !strings.HasPrefix(string(recvData), "You are connecting via") {

			return errors.New("服务器授权失败，未知错误")
		}

	}

	if !strings.Contains(string(recvData), "IPv4") {
		c.baseInfo.checkIp = 4
	} else {
		c.baseInfo.checkIp = 6
	}

	c.conn = conn

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
	if c.Interval == 0 {

		c.Interval = DefaultInterval
	}

	c.pingTime = make(map[string]uint)

	return nil
}
func (c *Client) getUpdateInfo() update {
	var ret update

	ret.Uptime = c.getUpTime()
	ret.Cpu = c.getCpuPercent()

	loadavg, err := load.Avg()
	if err != nil {
		loadavg.Load1 = 0
		loadavg.Load5 = 0
		loadavg.Load15 = 0
	}

	ret.Load1 = loadavg.Load1
	ret.Load5 = loadavg.Load5
	ret.Load15 = loadavg.Load15

	var memory = c.getMemory()
	ret.MemoryTotal = memory.total
	ret.MemoryUsed = memory.used

	var swap = c.getSwap()
	ret.SwapTotal = swap.total
	ret.SwapUsed = swap.used

	var hdd = c.getDiskUsage()
	ret.HddTotal = hdd.size
	ret.HddUsed = hdd.used

	var trafficData = c.getTraffic()
	ret.NetWorkIn = trafficData.in
	ret.NetWorkOut = trafficData.out

	ret.IpStatus = true

	var rateData = c.GetNetRate()
	ret.NetWorkRx = rateData.recvBytes / rateData.second
	ret.NetWorkTx = rateData.sendBytes / rateData.second

	ret.Ping10086 = c.getLostPacket("10086")
	ret.Ping10010 = c.getLostPacket("10010")
	ret.Ping189 = c.getLostPacket("189")

	ret.Time10086 = c.getPingTime("10086")
	ret.Time10010 = c.getPingTime("10010")
	ret.Time189 = c.getPingTime("189")

	var ioData = c.getDiskIo()
	ret.IoRead = ioData.readBytes / ioData.second
	ret.IoWrite = ioData.writeBytes / ioData.second

	var tupd = c.getTupd()
	ret.Tcp = tupd.tcp
	ret.Udp = tupd.udp
	ret.Process = tupd.process
	ret.Thread = tupd.thread

	c.lastUpdateTime = time.Now()

	return ret
}
