// client/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to listen for messages from the server
	go readMessages(conn)

	// Read user input and send messages to the server
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func readMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Print(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from server:", err)
	}
}
