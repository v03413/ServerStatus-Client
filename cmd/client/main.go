package main

import (
	"github.com/v03413/ServerStatus-Client/client"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	var debug bool
	var server, username, password, port string

	for _, v := range os.Args {
		if strings.HasPrefix(v, "SERVER=") {

			server = strings.TrimSpace(strings.Split(v, "SERVER=")[1])
		}
		if strings.HasPrefix(v, "PORT=") {

			port = strings.TrimSpace(strings.Split(v, "PORT=")[1])
		}
		if strings.HasPrefix(v, "USER=") {

			username = strings.TrimSpace(strings.Split(v, "USER=")[1])
		}
		if strings.HasPrefix(v, "PASSWORD=") {

			password = strings.TrimSpace(strings.Split(v, "PASSWORD=")[1])
		}
		if strings.HasPrefix(v, "DEBUG=") {

			debug = strings.TrimSpace(strings.Split(v, "DEBUG=")[1]) != ""
		}
	}

	c, err := client.NewClient(server, username, password, port, debug)
	if err != nil {

		log.Fatalln(err.Error())
	}

	go c.Start()

	log.Printf("开始运行，当前版本：%v", client.Version)

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
		runtime.GC()
	}
}
