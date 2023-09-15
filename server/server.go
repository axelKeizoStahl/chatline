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

type Profile struct {
    Name string
    Location string
    Description string
    Age int
}

type User struct {
    Serverlife bool
    Connection net.Conn
    Name string
    Profile Profile
    Room string
}

func (user *User) HandleMessage(c net.Conn, listener *Listener) {

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
            room := listener.User_rooms[message[12:len(message)-1]]
            room.Room_mutex.Lock()
            for index, element := range room.Users {
                if element.Connection == c {
                    room.Users[index] = room.Users[len(room.Users)-1]
                    room.Users = room.Users[:len(room.Users)-1]
                    room.Users = append(room.Users[:index], room.Users[index+1:]...)
                }
            }
            c.Close()
        }


        breaker, _ := regexp.Compile("83479256")
        breakpoints := breaker.FindAllStringSubmatchIndex(message, -1)
        user_name := message[:breakpoints[0][0]]
        room := listener.User_rooms[message[breakpoints[0][1]:breakpoints[1][0]]]
        fmt.Println(room, message[breakpoints[0][1]:breakpoints[1][0]])
        message = "[" + user_name + "] " + message[breakpoints[1][1]:]
        room.Room_mutex.Lock()
        defer room.Room_mutex.Unlock()
        room.Broadcast(message)

        room_assignment, _ := regexp.Compile("room_assign: ")
        if room_assignment.MatchString(message) {
            go func() {
                room_index := room_assignment.FindStringIndex(message)
                room_name := message[room_index[1]:len(message)-1]
                fmt.Println(room_name)
                user_name := message[:room_index[0]]
                if _, exists := listener.User_rooms[room_name]; !exists {
                    listener.User_rooms[room_name] = &Room{}
                }
                room := listener.User_rooms[room_name]
                room.Room_mutex.Lock()
                defer room.Room_mutex.Unlock()
                user := User{Name: user_name, Connection: c}
                room.Users = append(room.Users, user)
            }()
            continue
        }
    }
}

type Room struct {
    public bool
    Users []User
    Host User
    Room_mutex sync.Mutex
}

func (room *Room) HandleMessage(c net.Conn, listener *Listener) {

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
            room := listener.User_rooms[message[12:len(message)-1]]
            room.Room_mutex.Lock()
            for index, element := range room.Users {
                if element.Connection == c {
                    room.Users[index] = room.Users[len(room.Users)-1]
                    room.Users = room.Users[:len(room.Users)-1]
                    room.Users = append(room.Users[:index], room.Users[index+1:]...)
                }
            }
            c.Close()
        }


        breaker, _ := regexp.Compile("83479256")
        breakpoints := breaker.FindAllStringSubmatchIndex(message, -1)
        user_name := message[:breakpoints[0][0]]
        room := listener.User_rooms[message[breakpoints[0][1]:breakpoints[1][0]]]
        fmt.Println(room, message[breakpoints[0][1]:breakpoints[1][0]])
        message = "[" + user_name + "] " + message[breakpoints[1][1]:]
        room.Room_mutex.Lock()
        defer room.Room_mutex.Unlock()
        room.Broadcast(message)

        room_assignment, _ := regexp.Compile("room_assign: ")
        if room_assignment.MatchString(message) {
            go func() {
                room_index := room_assignment.FindStringIndex(message)
                room_name := message[room_index[1]:len(message)-1]
                fmt.Println(room_name)
                user_name := message[:room_index[0]]
                if _, exists := listener.User_rooms[room_name]; !exists {
                    listener.User_rooms[room_name] = &Room{}
                }
                room := listener.User_rooms[room_name]
                room.Room_mutex.Lock()
                defer room.Room_mutex.Unlock()
                user := User{Name: user_name, Connection: c}
                room.Users = append(room.Users, user)
            }()
            continue
        }
    }
}
func (r *Room) Broadcast(message string) {
    for _, conn := range r.Users {
        _, err := conn.Connection.Write([]byte(message))
        if err != nil {
            fmt.Println("Error broadcasting message:", err)
        }
    }
}

func (r *Room) ListUsers() {
    for _, conn := range r.Users {
        _, err := conn.Connection.Write([]byte(conn.Name + "\n"))
        if err != nil {
            fmt.Println("Error listing users:", err)
        }
    }
}

func ListRooms(rooms map[string]*Room) {
    for room_name := range rooms {
        fmt.Println(room_name)
    }
}

type Listener struct {
    Protocol string
    Port string
    User_rooms map[string]*Room
    Ln net.Listener
    Err error
    Exit_chan chan os.Signal
}

func (l *Listener) Listen() {
    l.Ln, l.Err = net.Listen(l.Protocol, l.Port)
    if l.Err != nil {
        fmt.Println(l.Err)
    }
}

func (l *Listener) Handle() {
    for {
        conn, err := l.Ln.Accept()
        if err != nil {
            fmt.Println(err)
        }
        go func(c net.Conn, listener *Listener) {
            defer c.Close()
            for {
                message, err := bufio.NewReader(c).ReadString('\n')
                if err != nil {
                    if err != io.EOF {
                        fmt.Println(err)
                    }
                    return
                }
                room_assignment, _ := regexp.Compile("room_assign: ")
                if room_assignment.MatchString(message) {
                    room_index := room_assignment.FindStringIndex(message)
                    room_name := message[room_index[1]:len(message)-1]
                    fmt.Println(room_name)
                    user_name := message[:room_index[0]]
                    if _, exists := l.User_rooms[room_name]; !exists {
                        l.User_rooms[room_name] = &Room{}
                    }
                    room := l.User_rooms[room_name]
                    room.Room_mutex.Lock()
                    defer room.Room_mutex.Unlock()
                    user := User{Name: user_name, Connection: c}
                    room.Users = append(room.Users, user)
                }
            }

        }(conn, l)
    }
}

func (l *Listener) Exit() {
    signal.Notify(l.Exit_chan, os.Interrupt)
    for range l.Exit_chan {
        fmt.Println("Closing all connections")
        for _, room := range l.User_rooms {
            go func(room *Room) {
                go func(room *Room) { room.Broadcast("Server is shutting down...\n") }(room)
                for _, user := range room.Users {
                    go func(user User) {
                        user.Connection.Close()
                        fmt.Println("Closed connection to", user.Name)
                    }(user)
                }
            }(room)
        }
    }
    fmt.Println("Closing server listener...")
    l.Ln.Close()
    fmt.Println("Exiting")
    os.Exit(1)
}

func Server() {
    fmt.Println("Starting server...")
    listener := Listener{Protocol: "tcp", Port: ":8000", User_rooms: make(map[string]*Room), Exit_chan: make(chan os.Signal, 1)}
    listener.Listen()

    go listener.Exit()
    for {
        conn, err := listener.Ln.Accept()
        if err!=nil {
            fmt.Println(err)
        }
        go func(c net.Conn) {
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
                    room := listener.User_rooms[message[12:len(message)-1]]
                    room.Room_mutex.Lock()
                    for index, element := range room.Users {
                        if element.Connection == c {
                            room.Users[index] = room.Users[len(room.Users)-1]
                            room.Users = room.Users[:len(room.Users)-1]
                            room.Users = append(room.Users[:index], room.Users[index+1:]...)
                        }
                    }
                    c.Close()
                }


                breaker, _ := regexp.Compile("83479256")
                breakpoints := breaker.FindAllStringSubmatchIndex(message, -1)
                user_name := message[:breakpoints[0][0]]
                room := listener.User_rooms[message[breakpoints[0][1]:breakpoints[1][0]]]
                fmt.Println(room, message[breakpoints[0][1]:breakpoints[1][0]])
                message = "[" + user_name + "] " + message[breakpoints[1][1]:]
                room.Room_mutex.Lock()
                defer room.Room_mutex.Unlock()
                room.Broadcast(message)
            }
        }(conn)
    }
}
