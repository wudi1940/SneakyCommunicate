package main

import (
	"cmd/serverCmd.go/internal"
	"flag"
)

var (
	serverIp   string
	serverPort int
)

// 通过命令行设置ip和端口号
// ./client -h，可以设置的默认信息
func init() {
	flag.StringVar(&serverIp, "ip", "0.0.0.0", "连接服务器IP（默认为127.0.0.1）")
	flag.IntVar(&serverPort, "port", 2968, "连接服务器Port（默认为2968）")
}

func main() {
	// 命令行解析
	flag.Parse()

	newServer := server.NewServer(serverIp, serverPort)
	newServer.StartServer()
}
