package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Serveraddr string
	Serverport int
	Clientname string
	Conn     net.Conn
	Flg      int
}

func Newclient(Serverip string,port int) *Client {
	newclient:=&Client{
		Serveraddr: Serverip,
		Serverport: port,
		Flg: 999,
		//Clientname: clientname,
	}
	conn,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",Serverip,port))
	if err!=nil{
		fmt.Println("连接错误",err)
		return nil
	}
	newclient.Conn=conn
	return newclient
}
//初始化
func init(){
	flag.StringVar(&Serverip,"Severip","127.0.0.1","输入目的服务器ip端口号，-Severip ：xxx.xxx.xxx")
	flag.IntVar(&Serverport,"severport",8888,"输入服务器端口号")
	//flag.StringVar(&clientname,"clientname","defaultValue","输入名字")

}
//全局变量
var Serverip string
var Serverport int
//var clientname string
func (client *Client)Ui()bool{
	var flg int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
    _,err:=fmt.Scanln(&flg)
	if err!=nil {
		fmt.Println("输入有误，请重新输入")
		return false
	}else if flg>=0&&flg<=3 {
		client.Flg=flg
		return true
	} else {
		fmt.Println("输入有误，请重新输入")
		return false
	}
}
func (client *Client)Run()  {
	for client.Flg!=0 {
		for client.Ui()!=true{}
		//根据不同模式处理不同业务
		switch client.Flg {
		case 1:
			client.Pubilc()

			fmt.Println("公聊模式选择")
		case 2:
			client.Tosomebody()
			fmt.Println("私聊模式选择")

		case 3:
			client.Rename()
			fmt.Println("您已经更改用户名为：",client.Clientname)

		}
	}

}
//处理Conn中的内容
func (client *Client)DealResponse(){
	io.Copy(os.Stdout,client.Conn)
}

//更改用户名
func (client *Client)Rename() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Clientname)
	sendMsg:="rename|"+client.Clientname+"\n"
	_,err:=	client.Conn.Write([]byte(sendMsg))
	if err!=nil {
		fmt.Println("重命名失败，清重新输入")
		return false
	}
	return true
}
func (client *Client) Showlist (){

}
func (client *Client) Tosomebody (){
	var somebody string
	var msg string
	fmt.Println(">>>>>>>您已进入私聊模式<<<<<<<")
	go client.DealResponse()


	fmt.Scanln(&somebody)
	for somebody!="c"{
		fmt.Println(">>>>>>>输入c退出私聊模式<<<<<<<")
		fmt.Scanln(&msg)
		for msg!="c" {
			if len(msg)!=0{
				_,err:=client.Conn.Write([]byte("to|"+somebody+"|"+msg+"\n"))
				if err!=nil{
					fmt.Println("输入有误，连接失败")
					break
				}
			}
			msg=""
			fmt.Scanln(&msg)
		}
		fmt.Scanln(&msg)
		fmt.Println(">>>>>>>请输入对方名字<<<<<<<")

	}

}
func (client *Client)Pubilc()  {
	var msg string
	fmt.Println(">>>>>>>您已进入公聊模式<<<<<<<")
	go client.DealResponse()
	fmt.Println(">>>>>>>输入c退出公聊模式<<<<<<<")
	fmt.Scanln(&msg)
	for msg!="c"{
		if len(msg)!=0 {
			_,err:=client.Conn.Write([]byte(msg+"\n"))
			if err !=nil{
				fmt.Println(">>>>>>>输入信息有误<<<<<<<")
				break
			}
		}
		msg=""
		//fmt.Println(":")
		fmt.Scanln(&msg)


		
	}

}
func main()  {
	flag.Parse()//解析命令行
	client:=Newclient("127.0.0.1",8888)
	if client==nil{
		fmt.Println("连接失败")
		return
	}else {
		fmt.Println("连接成功")
	}


	client.Run()
}
