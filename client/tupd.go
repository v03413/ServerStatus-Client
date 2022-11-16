package client

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

func (c *Client) getTupd(ret *update) {
	var r tupdStat

	if runtime.GOOS == "linux" {

		r = getOpenWrtTupd()
	} else {
		// windows
		r = getWinTupd()
	}

	ret.Tcp = r.tcp
	ret.Udp = r.udp
	ret.Process = r.process
	ret.Thread = r.thread
	c.waitGroup.Done()
}

func getOpenWrtTupd() tupdStat {
	var tcp, udp, process, thread uint64

	var validTcpSt = make(map[string]bool)
	validTcpSt["01"] = true
	validTcpSt["02"] = true
	validTcpSt["03"] = true
	validTcpSt["04"] = true
	validTcpSt["05"] = true
	validTcpSt["06"] = true
	validTcpSt["07"] = true
	validTcpSt["08"] = true
	//TCP_ESTABLISHED = 01
	//TCP_SYN_SENT    = 02
	//TCP_SYN_RECV    = 03
	//TCP_FIN_WAIT1   = 04
	//TCP_FIN_WAIT2   = 05
	//TCP_TIME_WAIT   = 06
	//TCP_CLOSE       = 07
	//TCP_CLOSE_WAIT  = 08
	//TCP_LAST_ACL = 09
	//TCP_LISTEN = 0A
	//TCP_CLOSING = 0B

	for _, file := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		if text, err := os.ReadFile(file); err == nil {
			lines := strings.Split(string(text), "\n")
			for _, line := range lines[1:] {
				fields := strings.Fields(line)
				if len(fields) > 4 {
					st := strings.TrimSpace(fields[3])
					if _, ok := validTcpSt[st]; ok {
						tcp += 1
					}
				}
			}
		}
	}

	for _, file := range []string{"/proc/net/udp", "/proc/net/udp6"} {
		if text, err := os.ReadFile(file); err == nil {
			lines := strings.Split(string(text), "\n")
			for _, line := range lines[1:] {
				fields := strings.Fields(line)
				if len(fields) > 4 {
					st := strings.TrimSpace(fields[3])
					if st != "07" { // udp socket 状态
						udp += 1
					}
				}
			}
		}
	}

	if hd, err := os.ReadDir("/proc"); err == nil {
		for _, file := range hd {
			if !isNumber(file.Name()) {

				continue
			}

			process += 1

			if text, err := os.ReadFile("/proc/" + file.Name() + "/stat"); err == nil {
				fields := strings.Fields(string(text))
				if t, err := strconv.ParseUint(fields[19], 10, 64); err == nil {

					thread += t
				}
			}
		}
	}

	return tupdStat{
		tcp:     uint(tcp),
		udp:     uint(udp),
		process: uint(process),
		thread:  uint(thread),
	}
}

func getWinTupd() tupdStat {

	return tupdStat{}
}

func isNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {

		return false
	}

	return true
}
