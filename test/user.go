package main

import (
	"fmt"
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

//创建一个用户的API
func NewUser(conn net.Conn,server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
		server: server,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User)  Offline(){
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name )
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "已xia线")
}
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
	//往当前的连接中直接发送字节，直接在终端中显示
}
//寻找某个用户 并将msg发给用户
func (this *User) Sendto(username string,msg string){
	remoteuer,ok:=this.server.OnlineMap[username]
	if !ok{
		this.SendMsg("当前用户不存在\n")
		return
	} else {
		remoteuer.SendMsg(msg+"\n")
		return
	}

}

//判断发送消息
func (this *User) Domessage(msg string){
	if msg=="who"{
		this.server.mapLock.Lock()
		for _,user:=range this.server.OnlineMap{
			onlinemsg:="["+user.Addr+"]"+":"+user.Name+"is online \n"
			this.SendMsg(onlinemsg)
		}

		this.server.mapLock.Unlock()
	}else if  len(msg)>len("RENAME|")&&msg[:7]=="rename|"{//重命名功能
		newname:=strings.Split(msg,"|")[1]
		//fmt.Println("newname:",newname)
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap,this.Name )
		this.server.OnlineMap[newname] = this
		this.server.mapLock.Unlock()
		this.Name=newname
		fmt.Println("maplock:",this.server.OnlineMap)
	} else if len(msg) >len("TO|")&&msg[:3]=="to|"{//私聊功能

		defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
			if err:=recover();err!=nil{
				this.SendMsg("输入错误,输入格式应该为to|username|message")
				return
			}
		}()
		username:=strings.Split(msg,"|")[1]
		usermsg:=strings.Split(msg,"|")[2]
		if usermsg=="" {
				this.SendMsg("message 不能为空")
			return
		}

		this.Sendto(username,usermsg)

	}else {
		this.server.BroadCast(this,msg)
	}

}
//监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
