package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
)

func Client() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
        for range c {
            fmt.Println()
            conn.Close()
            os.Exit(1)
        }
    }()
	writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(os.Stdin)
    netreader := bufio.NewReader(conn)

    fmt.Print("Enter room name: ")
    room, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    
    _, err = writer.WriteString("room_assign: " + room)
    if err != nil {
        fmt.Println("Error writing to buffered writer:", err)
        return
    }

    err = writer.Flush()
    if err != nil {
        fmt.Println("Error flushing buffered writer:", err)
        return
    }
     go func() {
        for {
            message, err := netreader.ReadString('\n')
            if err != nil {
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

        full_message := strconv.Itoa(len(room)-1) + room[:len(room)-1] + message

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

