// server/main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]string) // Track active clients and their names
	clientsMu sync.Mutex                  // Protect access to `clients`
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Chat server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Prompt for client name
	conn.Write([]byte("Please enter your name: "))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Println("Client disconnected before providing a name")
		return
	}

	clientName := scanner.Text()
	if clientName == "" {
		conn.Write([]byte("Invalid name. Connection closed.\n"))
		return
	}

	// Add client to the map with its name
	clientsMu.Lock()
	clients[conn] = clientName
	clientsMu.Unlock()

	fmt.Printf("Client connected: %s\n", clientName)
	broadcastMessage(fmt.Sprintf("Server: %s has joined the chat\n", clientName), conn)

	// Handle client messages
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Message from %s: %s\n", clientName, message)
		broadcastMessage(fmt.Sprintf("%s: %s\n", clientName, message), conn)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from %s: %v\n", clientName, err)
	}

	// Remove client and notify others
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()

	fmt.Printf("%s disconnected\n", clientName)
	broadcastMessage(fmt.Sprintf("Server: %s has left the chat\n", clientName), conn)
}

func broadcastMessage(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if client == sender {
			continue
		}
		client.Write([]byte(message))
	}
}
