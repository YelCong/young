package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8898")
	defer conn.Close()
	CheckErrorC(err)

	reader := bufio.NewReader(os.Stdin)
	buff := make([]byte, 1024)

	for {
		bytes, _, e := reader.ReadLine()
		CheckErrorC(e)

		conn.Write(bytes)
		fmt.Println("send message:" + string(bytes))

		//这段代码是先发再收，没发是走不到收的
		size, err := conn.Read(buff)
		CheckErrorC(err)

		msg := string(buff[0:size])
		fmt.Println("receive msg:" + msg)
	}

}

func CheckErrorC(err error) {
	if err != nil {
		fmt.Println("网络错误", err.Error())
		os.Exit(1)
	}
}
