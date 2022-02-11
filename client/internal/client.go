package client

import (
	"client/conf"
	"fmt"
	"net"

	"github.com/gen2brain/beeep"
)

type Client struct {
	ServerIp   string
	ServerPort int
	UserName   string
	flag       int
	Conn       net.Conn
}

// NewClient 新建一个客户端
func NewClient(serverIp string, port int) *Client {
	// 创建对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: port,
		flag:       999,
	}
	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, port))
	if err != nil {
		fmt.Println("net.dial err: ", err)
		return nil
	}

	client.Conn = conn
	return client
}

// Run 客户端运行逻辑
func (client *Client) Run() {
	for client.flag != conf.Exit {
		for client.menu() != true {
			fmt.Println("请输入对应的数字指令（0～4）")
		}
		// 根据不同模式处理不同的业务
		switch client.flag {
		case conf.PublicChat:
			// 群聊模式
			client.PublicChat()
		case conf.PrivateChat:
			// 私聊模式
			client.PrivateChat()
		case conf.UpdateName:
			// 更新用户名
			client.UpdateName()
		case conf.ShowOnlineUser:
			// 获取当前在线用户列表
			client.ShowOnlineUsers()
		}
	}
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.修改用户名")
	fmt.Println("4.查看在线用户列表")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	client.flag = flag

	return validate(flag)
}

func validate(flag int) bool {
	if flag >= 0 && flag <= 4 {
		return true
	} else {
		return false
	}
}

// PublicChat 公聊消息
func (client *Client) PublicChat() {
	var chatMsg string
	// 提示用户输入信息
	fmt.Printf(">>>>请输入聊天内容, 输入exit退出\n")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 消息不为空时发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Printf("conn.write error: %v", err)
				break
			}
		}
		chatMsg = ""
		fmt.Printf(">>>>请输入聊天内容, 输入exit退出\n")
		fmt.Scanln(&chatMsg)
	}
}

// PrivateChat 私聊
func (client *Client) PrivateChat() {
	// 选择私聊的用户名
	var remoteName string
	var chatMsg string
	// 展示在线用户
	client.ShowOnlineUsers()
	fmt.Printf(">>>>请输入用户名, 输入exit退出\n")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Printf(">>>>请输入消息内容, 输入exit退出\n")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空发送
			if len(remoteName) != 0 {
				sendMsg := "TO|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("PriviteChat.conn.write.error", err)
					break
				}
			}
			chatMsg = ""
			fmt.Printf(">>>>请输入消息内容, 输入exit退出\n")
			fmt.Scanln(&chatMsg)
		}
		client.ShowOnlineUsers()
		fmt.Printf(">>>>请输入用户名, 输入exit退出\n")
		fmt.Scanln(&remoteName)
	}
}

// UpdateName 修改用户名
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.UserName)

	sendMsg := "RENAME|" + client.UserName + "\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("con.Write error: ", err)
		return false
	}
	return true
}

// ShowOnlineUsers 获取在线用户列表
func (client *Client) ShowOnlineUsers() {
	sendMsg := "WHO\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Printf("ShowOnlineUsers conn.write error: %v", err)
		return
	}
}

// ListenMsg  监听服务器返回的消息
func (client *Client) ListenMsg() {
	// 一旦client.conn有数据，就直接copy到stout标准输出
	// 永久阻塞，等同于下面for循环
	//io.Copy(os.Stdout, client.Conn)
	for {
		buf := make([]byte, 1024)
		n, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}

		msg := string(buf[:n])

		err = beeep.Notify("你有新的消息", msg, "assets/information.png")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(msg)
	}
}
