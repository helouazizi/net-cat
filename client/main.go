// client/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("You are Connected to the server")

	// Read server messages in a separate goroutine
	go func() {
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Println(serverScanner.Text())
		}
		if err := serverScanner.Err(); err != nil {
			fmt.Println("Error reading from server:", err)
		}
	}()

	// Send user input to the server
	//fmt.Println("Your :")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		message := scanner.Text()
		conn.Write([]byte(message + "\n"))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
