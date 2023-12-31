package client

import (
    "io"
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func help() {
    fmt.Println("use ** as a prefix to enter a command")
    fmt.Println("**logo : log out")
    fmt.Println("**exit : exit chatline")
    fmt.Println("**switchu : switch user")
}

func Client() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

    fmt.Println("WELCOME to CHATLINE")
    fmt.Println("type **h to ask for help")
	writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(os.Stdin)
    netreader := bufio.NewReader(conn)

    var username string
    fmt.Print("Enter username: ")
    fmt.Scanf("%s", &username)

    fmt.Print("Enter room name: ")
    room, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    
    _, err = writer.WriteString(username + "room_assign: " + room)
    if err != nil {
        fmt.Println("Error writing to buffered writer:", err)
        return
    }

    err = writer.Flush()
    if err != nil {
        fmt.Println("Error flushing buffered writer:", err)
        return
    }

    exit := func() {
        _, err = writer.WriteString("83479256exit" + room)
        if err != nil {
            fmt.Println("Error writing to buffered writer:", err)
            return
        }
        fmt.Println()
        conn.Close()
        os.Exit(1)
    }

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
        for range c {
            exit()
        }
    }()

     go func() {
        for {
            message, err := netreader.ReadString('\n')
            if err != nil {
                if err == io.EOF {
                    exit()
                }
                fmt.Println("Error reading message:", err)
                return
            }
            fmt.Print("\nReceived: " + message)
            fmt.Print("Enter a message: ")
        }
    }()
	for {
		fmt.Print("Enter a message: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

        full_message := username + "83479256" + room[:len(room)-1] + "83479256" + message
        if message[:len(message)-1] == "**h" {
            help()
        }
        if message[:len(message)-1] == "**exit" {
            exit()
        }

		_, err = writer.WriteString(full_message)
		if err != nil {
			fmt.Println("Error writing to buffered writer:", err)
			return
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing buffered writer:", err)
			return
		}
	}
}

