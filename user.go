package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string

	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	user.server.BroadCast(user, "online")
}

func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	user.server.BroadCast(user, "offline")
}

func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

func (user *User) DoMsg(msg string) {
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, val := range user.server.OnlineMap {
			onlineMsg := "[" + val.Addr + "]" + val.Name + ":" + "is online \n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := user.server.OnlineMap[newName]

		if ok {
			user.SendMsg(newName + "is existing\n")
			return
		}
		user.server.mapLock.Lock()
		delete(user.server.OnlineMap, user.Name)
		user.server.OnlineMap[newName] = user
		user.Name = newName
		user.server.mapLock.Unlock()

		user.SendMsg("username : [" + newName + "] has been updated\n")

	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("personal communication msg format is wrong, please process msg like that 'to|username|msg'\n")
			return
		}

		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok{
			user.SendMsg(remoteName + " is not online or not exist")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			user.SendMsg("your msg is empty, why not write something to your friend\n")
			return
		}

		remoteUser.SendMsg(user.Name + " msg : " + content)

	} else {
		user.server.BroadCast(user, msg)

	}
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
