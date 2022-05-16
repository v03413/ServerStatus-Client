package client

import (
	"github.com/shirou/gopsutil/host"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func (c *Client) getTupd() tupdStat {
	info, err := host.Info()
	if err != nil {
		log.Println(err.Error())

		return tupdStat{}
	}

	if info.Platform == "openwrt" {

		return getOpenWrtTupd()
	}

	if strings.HasPrefix(runtime.GOOS, "linux") {

		return getLinuxTupd()
	}

	if strings.HasPrefix(runtime.GOOS, "win") {

		return getWinTupd()
	}

	return tupdStat{}
}

func getOpenWrtTupd() tupdStat {
	var tcp, udp, process, thread uint64

	cmd := exec.Command("bash", "-c", "netstat -t | wc -l")
	if out, err := cmd.CombinedOutput(); err != nil {

		log.Println(err.Error())
	} else {
		tcp, _ = strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	}

	cmd = exec.Command("bash", "-c", "netstat -u | wc -l")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println(err.Error())
	} else {
		udp, _ = strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	}

	if hd, err := ioutil.ReadDir("/proc"); err == nil {
		for _, file := range hd {
			if !isNumber(file.Name()) {

				continue
			}

			process += 1

			if text, err := ioutil.ReadFile("/proc/" + file.Name() + "/stat"); err == nil {
				fields := strings.Fields(string(text))
				if t, err := strconv.ParseUint(fields[19], 10, 64); err == nil {

					thread += t
				}
			}
		}
	}

	return tupdStat{
		tcp:     uint(tcp - 2),
		udp:     uint(udp - 2),
		process: uint(process),
		thread:  uint(thread),
	}
}
func getLinuxTupd() tupdStat {
	// 这里偷懒了，有空再说吧
	return getOpenWrtTupd()

	//return tupdStat{}
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
