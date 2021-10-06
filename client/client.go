package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type ServerConfig struct {
	Ip      string // server ip
	Port    int    // server port
	Network string // network tcp/udp
}

func send(conn net.Conn, ch chan int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("reader input failed")
		}
		data = strings.Trim(data, "\r\n")
		data = strings.Trim(data, " ")
		if data == "exit" || data == "quit" {
			ch <- 1
			return
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("send message failed")
		}
	}
}

func receive(conn net.Conn) {
	for {
		buff := make([]byte, 1024)
		len, err := conn.Read(buff)
		if err != nil {
			fmt.Println("receive message failed")
		}
		fmt.Println(string(buff[:len]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5678")
	if err != nil {
		fmt.Println("connect to server failed")
		return
	}
	defer conn.Close()
	ch := make(chan int)
	go send(conn, ch)
	go receive(conn)
	for {
		if (<-ch) == 1 {
			fmt.Println("exited!")
			return
		}
	}
}
