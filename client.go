package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

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

