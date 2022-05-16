package main

import (
	client2 "ServerStatus-Client/client"
	"fmt"
)

func main() {
	var server = "192.168.2.120"
	var port = 5566
	var username = "r2s"
	var password = "123234zxc"
	var interval = 1
	var client = client2.Client{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
		Interval: uint(interval),
	}

	err := client.Start()
	if err != nil {

		fmt.Println(err.Error())
	}
}
