package main

import (
	"fmt"
	"net"
)

var (
	Ip      string
	Port    int
	Clients map[string]net.Conn
)

func init() {
	Ip = ""
	Port = 5678
	Clients = make(map[string]net.Conn)
}

func broadCast(clients map[string]net.Conn, message string, operator net.Conn) {
	fmt.Println(operator.RemoteAddr().String(), " boardcast message:", message)
	for _, client := range clients {
		if client == operator {
			continue
		}
		_, err := client.Write([]byte(message))
		if err != nil {
			fmt.Println("message send to ", client.RemoteAddr().String(), " failed!")
		}
	}
}

func handConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("client %s connected\n", conn.RemoteAddr().String())
	for {
		buff := make([]byte, 1024)
		len, err := conn.Read(buff)
		if err != nil {
			delete(Clients, conn.RemoteAddr().String())
			fmt.Println(conn.RemoteAddr().String(), " off line")
			return
		}
		// broadcast
		broadCast(Clients, string(buff[:len]), conn)
	}
}

func main() {
	// create socket
	listen, err := net.Listen("tcp", ":5678")
	if err != nil {
		fmt.Println("server start failed")
		return
	}
	defer listen.Close()
	fmt.Println("server started!")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("client connect failed failed")
			return
		}
		Clients[conn.RemoteAddr().String()] = conn
		go handConnection(conn)
	}
}
