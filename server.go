package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	server.BroadCast(user, "user online")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				server.BroadCast(user, "user offline")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err : ", err)
				return
			}

			msg := string(buf[:n-1])
			server.BroadCast(user, msg)
		}
	}()
}

func (server *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err : ", err)
		return
	}
	defer listen.Close()

	go server.ListenMessage()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err : ", err)
			continue
		}

		go server.Handler(conn)
	}
}
