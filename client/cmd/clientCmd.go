package main

import (
	"client/internal"
	"flag"
	"fmt"
)

var (
	serverIp   string
	serverPort int
)

// 通过命令行设置ip和端口号
// ./client -h，可以设置的默认信息
func init() {
	flag.StringVar(&serverIp, "ip", "101.43.177.135", "连接服务器IP（默认为101.43.177.135）")
	flag.IntVar(&serverPort, "port", 2968, "连接服务器Port（默认为2968）")
}

func main() {
	// 命令行解析
	flag.Parse()

	newClient := client.NewClient(serverIp, serverPort)
	if newClient == nil {
		fmt.Println(">>>>>>>>连接服务器失败")
		return
	}
	fmt.Println(">>>>>>>>连接成功")

	go newClient.ListenMsg()
	// 启动客户端业务
	newClient.Run()
}
