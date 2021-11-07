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

	go user.ListenMessgae()

	return user
}

func (this *User) Online() {
	// 用户上线
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播当前用户上线消息
	this.server.BoardCast(this, "已上线")
}

func (this *User) Offline() {
	// 用户下线
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 广播当前用户下线消息
	this.server.BoardCast(this, "已下线")
}

func (this *User) SendMessgae(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessgae(msg string) {
	if msg == "who" {
		// 查询当前在线人数
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "{" + user.Addr + "}" + user.Name + ":在线...\n"
			this.SendMessgae(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名
		newName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessgae("当前用户名已占用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessgae("您已用户名" + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊
		target := strings.Split(msg, "|")[1] // 获取目标用户名
		user, ok := this.server.OnlineMap[target]
		if ok {
			message := strings.Split(msg, "|")[2]
			user.SendMessgae(this.Name + "对您说: " + message)
		}
	} else {
		this.server.BoardCast(this, msg)
	}
}

// 监听当前 channel 的方法， 一旦有消息， 就直接发送给客户端
func (this *User) ListenMessgae() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
