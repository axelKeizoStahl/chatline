package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
    "os/signal"
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

	for {
		fmt.Print("Enter a message: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		_, err = writer.WriteString(message)
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

