package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8898")
	defer listener.Close()
	CheckErrorS(err)
	for {
		conn, err := listener.Accept() //新的客户端连接
		CheckErrorS(err)
		go ProcessInfo(conn)
	}

}

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 10)
	defer conn.Close()

	for {
		size, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Println(time.Now())

		if size != 0 {
			msg := string(buf[0:size])
			fmt.Println("收到消息:" + msg)
			_, _ = conn.Write([]byte("已阅.."))
			if msg == "break" {
				break
			}
		}
	}
	fmt.Println("ending")
	return
}

func CheckErrorS(err error) {
	if err != nil {
		fmt.Println("网络错误", err.Error())
		os.Exit(1)
	}
}
