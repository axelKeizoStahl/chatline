package server

import (
    "strconv"
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

func (r *Room) Broadcast(message string) {
    for _, conn := range r.Connections {
        _, err := conn.Write([]byte(message))
        if err != nil {
            fmt.Println("Error broadcasting message:", err)
        }
    }
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

                is_exit, _ := regexp.Compile("83479256exit")
                if is_exit.MatchString(message) {
                    room_mutex.Lock()
                    room := user_rooms[string(message)[12:len(string(message))]]
                    for index, element := range room.Connections {
                        if element == c {
                            room.Connections[index] = room.Connections[len(room.Connections)-1]
                            room.Connections = room.Connections[:len(room.Connections)-1]
                            room.Connections = append(room.Connections[:index], room.Connections[index+1:]...)
                        }
                    }
                    c.Close()
                }


                room_assignment, _ := regexp.Compile("room_assign: .*")
                if room_assignment.MatchString(message) {
                    go func() {
                        room_mutex.Lock()
                        defer room_mutex.Unlock()
                        if _, exists := user_rooms[message[13:len(message)-1]]; !exists {
                            user_rooms[message[13:len(message)-1]] = &Room{}
                        }
                        room := user_rooms[message[13:len(message)-1]]
                        room.Connections = append(user_rooms[message[13:len(message)-1]].Connections, c)
                    }()
                    continue
                }
                num, _ := strconv.Atoi(string(message)[:1])
                fmt.Print("Message Received: ", string(message)[num+1:])
                room := user_rooms[string(message)[:num+1][1:]]
                room_mutex.Lock()
                room.Broadcast(string(message[num+1:]))
                room_mutex.Unlock()
            }
        }(conn)
    }
}
