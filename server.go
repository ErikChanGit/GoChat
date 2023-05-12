package main

//
import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 编译： go build -o server a_sample.go a_server.go
// 启动服务： ./server
// 连接服务器： nc 127.0.0.1 8888

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) ListMessage() {
	for {
		// 从通道获取消息
		msg := <-this.Message
		this.mapLock.Lock()
		fmt.Println(this.OnlineMap)
		for _, user := range this.OnlineMap {
			fmt.Println(user.Name + "  " + msg)
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BoardCast(user *User, msg string) {
	sendMsg := "{" + user.Addr + "}" + user.Name + ":" + msg
	// fmt.Println(sendMsg)
	// 发送消息到通道
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("连接建立成功")
	isLive := make(chan bool)

	user := NewUser(conn, this)

	user.Online()

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for { // 无限循环等待消息
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				this.BoardCast(user, "已下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessgae(msg)

			isLive <- true
		}
	}()

	for {
		select {
		// 强制剔除
		case <-isLive:

		case <-time.After(time.Second * 10000): // 重置定时器
			user.Offline()
			user.SendMessgae("10000s 强制下线")
			close(user.C) // 关闭管道
			conn.Close()  // 关闭连接
			return
		} // 阻塞， 以免退出此函数
	}

}

func (this *Server) Start() {
	// socket 监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer listener.Close()

	go this.ListMessage()
	// accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

	// close listen socket
}
