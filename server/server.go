package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"regexp"
	"sync"
)

type Room struct {
    Connections []net.Conn
}

func Server() {
    fmt.Println("Starting server...")
    listener, err := net.Listen("tcp", "127.0.0.1:8000")
    if err!=nil {
        fmt.Println(err)
    }
    defer func() { _ = listener.Close() }()
    user_rooms := make(map[string]*Room)
    var room_mutex sync.Mutex

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
                message, err := bufio.NewReader(c).ReadString('\n')
                if err != nil {
                    if err != io.EOF {
                        fmt.Println(err)
                    }
                    return
                }

                room_assignment, _ := regexp.Compile("room_assign: .*")
                if room_assignment.MatchString(message) {
                    go func() {
                        room_mutex.Lock()
                        defer room_mutex.Unlock()
                        if _, exists := user_rooms[message[14:]]; !exists {
                            user_rooms[message[14:]] = &Room{}
                        }
                        room := user_rooms[message[14:]]
                        room.Connections = append(user_rooms[message[14:]].Connections, c)
                    }()
                    continue
                }
                fmt.Println("Message Received: ", string(message))
            }
        }(conn)
    }
}
