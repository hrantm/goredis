package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		processRequest(conn)

	}

}

func processRequest(conn net.Conn) {
	defer conn.Close()

	resp := newResp(conn)

	value, err := resp.Read()
	if err != nil {
		fmt.Println("Error: ", err)
		message := fmt.Sprintf("-%v\r\n", err)
		conn.Write([]byte(message))
		return
	}
	fmt.Println(value)
	conn.Write([]byte("+OK\r\n"))

}
