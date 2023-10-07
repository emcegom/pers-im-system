package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : " , err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip, and the default value is 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "set server port, and the default value is 8888")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println("connecting server " + serverIp + ":" + strconv.Itoa(serverPort) + " fail")
		return
	}
	fmt.Println("connecting server " + serverIp + ":" + strconv.Itoa(serverPort) + " success")
	select{
	}
}
