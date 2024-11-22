// server/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients = make(map[net.Conn]bool) // To keep track of connected clients
	mu      sync.Mutex                // Mutex to protect the clients map
)

func main() {
	// Start our TCP server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server running on port 8080..")

	for {
		// Accept incoming connections
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		clients[connection] = true // Add new client to the map
		mu.Unlock()

		go handleClient(connection)
	}
}

func handleClient(connection net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, connection) // Remove client from the map on disconnect
		mu.Unlock()
		connection.Close()
	}()

	// Prompt the client to enter their name
	connection.Write([]byte("Enter your name: "))
	reader := bufio.NewReader(connection)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	name = name[:len(name)-1] // Remove the newline character
	fmt.Printf("%s joined the server\n", name)

	broadcastMessage(fmt.Sprintf("%s has joined the chat!\n", name))

	for {
		// Prompt the client to enter their message
		connection.Write([]byte("Enter your message: "))
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		message = message[:len(message)-1] // Remove the newline character
		broadcastMessage(fmt.Sprintf("%s: %s\n", name, message))
	}
}

func broadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		_, err := client.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message to client:", err)
		}
	}
}
