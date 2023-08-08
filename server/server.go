package server

import (
	"bufio"
	"fmt"
	"net"
)

func Server() {
    fmt.Println("Starting server...")
    ln, _ := net.Listen("tcp", ":8000")
    conn, _ := ln.Accept()
    defer conn.Close()

    for {
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Println("Message Received: ", string(message))
    }
}
