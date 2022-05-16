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

	for _, v := range os.Args {
		if strings.HasPrefix(v, "SERVER=") {

			client.Server = strings.TrimSpace(strings.Split(v, "SERVER=")[1])
		}
		if strings.HasPrefix(v, "PORT=") {

			client.Port, _ = strconv.ParseUint(strings.TrimSpace(strings.Split(v, "PORT=")[1]), 10, 64)
		}
		if strings.HasPrefix(v, "USERNAME=") {

			client.Username = strings.TrimSpace(strings.Split(v, "USERNAME=")[1])
		}
		if strings.HasPrefix(v, "PASSWORD=") {

			client.Password = strings.TrimSpace(strings.Split(v, "PASSWORD=")[1])
		}
	}

	err := client.Start()
	if err != nil {

		log.Println(err.Error())
	}

	runtime.GC()
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
	}
}
