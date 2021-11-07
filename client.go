package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// 客户端结构体
type Client struct {
	Ip   string
	Port int
	Name string
	Conn net.Conn
	flag int
}

// 新建客户端
func NewClient(ip string, port int) *Client {
	client := &Client{ // 构造客户端
		Ip:   ip,
		Port: port,
		flag: 999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port)) // 创建 tcp 连接
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.Conn = conn

	return client
}

/**
主菜单
*/
func (client *Client) menu() bool {

	var _flag int

	fmt.Println("-----腾讯会议------")
	fmt.Println("1、公聊模式")
	fmt.Println("2、私聊模式")
	fmt.Println("3、更新用户名")

	fmt.Scanln(&_flag)

	if _flag >= 0 && _flag <= 3 {
		client.flag = _flag
		return true
	} else {
		fmt.Println("")
		return false
	}
}

// 更新用户名
func (client *Client) updateName() bool {

	fmt.Println(">>输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.Conn.Write([]byte(sendMsg))

	if err != nil {
		return false
	}

	client.menu()
	return true
}

// 公聊模式
func (client *Client) publicCHat() {
	var msg string
	fmt.Println(">>输入内容:")
	fmt.Scanln(&msg)

	for msg != "exit" {
		sendMsg := msg + "\n"
		_, err := client.Conn.Write([]byte(sendMsg))
		if err != nil {
			break
		}

		msg = ""
		fmt.Println(">>输入内容:")
		fmt.Scanln(&msg)
	}

	client.menu()
}

// 私聊模式
func (client *Client) privateChat() {

	sendMsg := "who\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		return
	}

	var user string
	fmt.Println(">>输入聊天对象:")
	fmt.Scanln(&user)

	for user != "exit" {
		var msg string
		fmt.Println(">>输入聊天消息:")
		fmt.Scanln(&msg)

		for msg != "exit" {
			sendMsg := "to|" + user + "|" + msg + "\n"
			_, err := client.Conn.Write([]byte(sendMsg))
			if err != nil {
				break
			}

			msg = ""
			fmt.Println(">>继续输入内容:")
			fmt.Scanln(&msg)
		}

		// 循环提问
		sendMsg = "who\n"
		_, err = client.Conn.Write([]byte(sendMsg))
		if err != nil {
			return
		}

		var user string
		fmt.Println(">>输入聊天对象:")
		fmt.Scanln(&user)

	}

	client.menu()
}

// 处理回应
func (client *Client) DealResponse() {
	// 将客户端的消息输出
	// 在文件指针之间直接复制的，不用全读入内存
	io.Copy(os.Stdout, client.Conn)
}

// Client 结构体对象的函数
func (client *Client) run() {
	if client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {
		case 1:
			fmt.Println("选中公聊模式")
			client.publicCHat()
			break
		case 2:
			fmt.Println("选中私聊模式")
			client.privateChat()
			break
		case 3:
			fmt.Println("更新用户名")
			client.updateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// go 文件首次自动进入的函数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置 IP 地址")
	flag.IntVar(&serverPort, "port", 8888, "设置 端口")
}

func main() {
	flag.Parse() // 把用户传递的命令行参数解析为对应变量的值

	client := NewClient(serverIp, serverPort)
	// fmt.Println(client)
	if client == nil {
		fmt.Println("服务器连接失败")
		return
	}

	// go 关键字并发执行
	go client.DealResponse()
	fmt.Println("服务器连接成功")

	client.run()
	select {}
}
