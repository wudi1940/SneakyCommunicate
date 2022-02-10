package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	path string
	port int
)

// 通过命令行设置ip和端口号
// ./client -h，可以设置的默认信息
func init() {
	flag.StringVar(&path, "path", "./", "文件路径，默认./")
	flag.IntVar(&port, "port", 80, "服务器Port（默认为80）")
}

func main() {
	flag.Parse()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Dir not exist")
	}

	http.Handle("/", http.FileServer(http.Dir(path)))

	s := &http.Server{
		Addr:        fmt.Sprintf(":%v", port),
		ReadTimeout: 150 * time.Second,
	}

	fmt.Printf("Start listen on port %d\n", port)

	log.Fatal(s.ListenAndServe())
}
