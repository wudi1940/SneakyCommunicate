package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int

	// 在线用户列表， 当有新用户连接时， 将用户放入map中
	// 由于map是线程不安全类， 需要加锁
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 接受来自用户的消息
	MessageChannel chan string
}

// NewServer 初始化一个服务器对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:             ip,
		Port:           port,
		OnlineMap:      make(map[string]*User),
		MessageChannel: make(chan string),
	}
	return server
}

// StartServer 启动服务器对象
func (server *Server) StartServer() {
	fmt.Printf("服务器启动, IP:%s port:%d", server.IP, server.Port)
	// socket listen
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("net.listen err: ", err)
		return
	}
	defer listener.Close()

	// 消息监听
	go server.ListenMsg()

	// 有客户端接入时的处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}
		go server.ClientHandle(conn)
	}

}

// ClientHandle 客户端接入处理
func (server *Server) ClientHandle(conn net.Conn) {
	// 当前连接业务
	fmt.Println("连接建立成功")
	// 为连接创建新用户
	newUser := NewUser(conn, server)
	// 新用户接入后 创建User放入map
	newUser.Register()
	// 监听用户活跃的channel
	isLive := make(chan bool)

	// 开启一个gouroutine，接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		// 监听
		for {
			n, err := conn.Read(buf)
			// 报错
			if err != nil && err != io.EOF {
				fmt.Println("con read Error", err)
				return
			}
			// 0时，表示用户下线
			if n == 0 {
				newUser.UnRegister()
				return
			}
			// 接受消息，去除结尾/n
			msg := string(buf[:n-1])
			// 发送msg
			newUser.MsgHandle(msg)

			// 用户发送消息，代表其处于活跃状态
			isLive <- true
		}
	}()

	// 阻塞当前handler，否则当前goroutine会关闭，里面的子goroutine也关闭
	for {
		select {
		case <-isLive:
			// 当前用户活跃，重置定时器
			// 不在任何事，为了激活select，更新下面的定时器
		case <-time.After(time.Minute * 60):
			// 已经过时
			// 下线用户
			newUser.SendMsg("长时间afk，被下线，请重新登陆")
			newUser.UnRegister()
			// 关闭连接
			conn.Close()
			// 退出当前handler
			return
			// or runtime.Goexit()
		}
	}
}

// Broadcast 广播消息给所有用户
func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + msg
	//fmt.Println("<<<<<<<<<<" + msg)
	server.MessageChannel <- sendMsg
}

// ListenMsg 监听message，一旦有消息就发送到所有user，具体消息处理逻辑在user.MsgHandle中
func (server *Server) ListenMsg() {
	for {
		msg := <-server.MessageChannel
		// msg发送给所有在线用户
		server.mapLock.Lock()
		for _, client := range server.OnlineMap {
			client.MsgChan <- msg
		}
		server.mapLock.Unlock()
	}
}
