package server

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	conn    net.Conn

	// server
	Server *Server
}

// NewUser 创建新用户对象
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Addr:    userAddr,
		MsgChan: make(chan string),
		conn:    conn,
		Server:  server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

// Register 用户上线注册
func (user *User) Register() {

	user.Server.mapLock.Lock()
	user.Server.OnlineMap[user.Name] = user
	user.Server.mapLock.Unlock()

	// 广播上线提醒
	user.Server.Broadcast(user, "用户上线")

}

// UnRegister 用户下线注销
func (user *User) UnRegister() {

	user.Server.mapLock.Lock()
	delete(user.Server.OnlineMap, user.Name)
	user.Server.mapLock.Unlock()

	close(user.MsgChan)
	// 广播下线提醒
	user.Server.Broadcast(user, "用户离线")
}

// MsgHandle 处理用户消息
func (user *User) MsgHandle(msg string) {
	// 查询在线用户，并返回给当前用户
	if msg == "who" {
		user.GetOnlineUsers()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		user.Rename(msg)
	} else if len(msg) > 3 && msg[:3] == "to|" {
		user.PrivateChat(msg)
	} else {
		user.Server.Broadcast(user, msg)
	}
}

// PrivateChat 私聊
func (user *User) PrivateChat(msg string) {
	// 获取to用户名
	toUserName := strings.Split(msg, "|")[1]
	if toUserName == "" {
		user.SendMsg("消息格式不正确，请使用\"to|张三|消息\"的格式")
		return
	}
	toUser, ok := user.Server.OnlineMap[toUserName]
	if !ok {
		user.SendMsg("该用户不存" + "\n")
		return
	}
	toMsg := strings.Split(msg, "|")[2]
	if toMsg == "" {
		user.SendMsg("请输入消息内容" + "\n")
		return
	}
	toUser.SendMsg(toMsg)
}

// GetOnlineUsers 获取当前在线用户列表
func (user *User) GetOnlineUsers() {
	user.Server.mapLock.Lock()
	for _, client := range user.Server.OnlineMap {
		user.MsgChan <- "[" + client.Addr + "] " + client.Name + ":online..."
	}
	user.Server.mapLock.Unlock()
}

// Rename 修改用户名
func (user *User) Rename(newName string) {
	// 检查是否被占用
	newName = newName[7:]
	_, ok := user.Server.OnlineMap[newName]
	if ok {
		user.SendMsg("【" + newName + "】" + "已被占用" + "\n")
	} else {
		user.Server.mapLock.Lock()
		delete(user.Server.OnlineMap, user.Name)
		user.Server.OnlineMap[newName] = user
		user.Server.mapLock.Unlock()
		user.Name = newName
		user.SendMsg("修改成功，当前用户名：" + newName + "\n")
	}
}

// SendMsg 对当前用户发送消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// ListenMessage 监听当前User channel的方法，有消息就发给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.MsgChan
		user.SendMsg(msg + "\n")
	}
}
