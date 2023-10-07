package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
		fmt.Println("net.Dial err : ", err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println(
		`
		1.公聊模式
		2.私聊模式
		3.更新用户名
		0.退出
		`,
	)

	fmt.Scan(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please input number in the range")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		}
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("please input a new username : ")

	var newName string

	fmt.Scanln(&newName)

	sendMsg := "rename|" + newName + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	client.Name = newName
	return true

}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("please input your message or input exit to logout")

	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {

		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err : ", err)
				break
			}

		}

		chatMsg = ""
		fmt.Println("please input your message or input exit to logout")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err : ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string

	client.SelectUsers()

	fmt.Println("please input the communication part, or input exit to logout")

	fmt.Scanln(&remoteName)

	for remoteName != "exit" {

		var chatMsg string

		fmt.Println("please input your message or input exit to logout")

		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {

			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err : ", err)
					break
				}

			}

			chatMsg = ""
			fmt.Println("please input your message or input exit to logout")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()

		fmt.Println("please input the communication part, or input exit to logout")

		fmt.Scanln(&remoteName)

	}

}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
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

	go client.DealResponse()

	fmt.Println("connecting server " + serverIp + ":" + strconv.Itoa(serverPort) + " success")

	client.Run()
}
