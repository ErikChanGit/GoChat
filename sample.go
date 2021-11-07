package main

// import "fmt"

// 服务端

// OnlineMap

// Channel

// 客户端

//

const ( // ( 不是 {
	BJ = iota // 值为0， 后面会自增
	SH
	SZ
)

func main() {
	// const length int = 10
	// fmt.Println(SZ)
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
