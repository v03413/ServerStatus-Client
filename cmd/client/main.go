package main

import (
	client2 "ServerStatus-Client/client"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	var client = client2.Client{}
	var osSignals = make(chan os.Signal, 1)

	for _, v := range os.Args {
		if strings.HasPrefix(v, "SERVER=") {

			client.Server = strings.TrimSpace(strings.Split(v, "SERVER=")[1])
		}
		if strings.HasPrefix(v, "PORT=") {

			client.Port, _ = strconv.ParseUint(strings.TrimSpace(strings.Split(v, "PORT=")[1]), 10, 64)
		}
		if strings.HasPrefix(v, "USER=") {

			client.Username = strings.TrimSpace(strings.Split(v, "USER=")[1])
		}
		if strings.HasPrefix(v, "PASSWORD=") {

			client.Password = strings.TrimSpace(strings.Split(v, "PASSWORD=")[1])
		}
		if strings.HasPrefix(v, "DEBUG=") {

			client.Debug = strings.TrimSpace(strings.Split(v, "DEBUG=")[1]) != ""
		}
	}

	log.Printf("开始运行，当前版本：%v", client2.Version)

	err := client.Start()
	if err != nil {

		log.Println(err.Error())
		signal.Notify(osSignals, os.Interrupt)
		return
	}

	runtime.GC()
	{
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
	}
}
