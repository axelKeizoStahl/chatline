package server

import (
    "io"
	"bufio"
	"fmt"
	"net"
    "os"
    "os/signal"
)

func Server() {
    fmt.Println("Starting server...")
    listener, err := net.Listen("tcp", "127.0.0.1:8000")
    if err!=nil {
        fmt.Println(err)
    }
    defer func() { _ = listener.Close() }()

    exit := make(chan os.Signal, 1)
    signal.Notify(exit, os.Interrupt)
    go func() {
        for range exit {
            fmt.Println()
            listener.Close()
            os.Exit(1)
        }
    }()

    for {
        conn, err := listener.Accept()
        if err!=nil {
            fmt.Println(err)
        }
        go func(c net.Conn) {
            exit := make(chan os.Signal, 1)
            signal.Notify(exit, os.Interrupt)
            go func() {
                for range exit {
                    fmt.Println()
                    c.Close()
                    os.Exit(1)
                }
            }()
            defer c.Close()
            for {
                message, err := bufio.NewReader(conn).ReadString('\n')
                if err != nil {
                    if err != io.EOF {
                        fmt.Println(err)
                    }
                    return
                }
                fmt.Println("Message Received: ", string(message))
            }
        }(conn)
    }
}
