package server

import (
	"bufio"
	"fmt"
	"net"
)

func Server() {
    fmt.Println("Starting server...")
    listener, err := net.Listen("tcp", "127.0.0.1:8000")
    if err!=nil {
        fmt.Println(err)
    }
    defer func() { _ = listener.Close() }()

    for {
        conn, err := listener.Accept()
        if err!=nil {
            fmt.Println(err)
        }
        go func(c net.Conn) {
            defer c.Close()
            for {
                message, _ := bufio.NewReader(conn).ReadString('\n')
                if message != "" {
                    fmt.Println("Message Received: ", string(message))
                } else {
                    break
                }
            }
        }(conn)
    }
}
