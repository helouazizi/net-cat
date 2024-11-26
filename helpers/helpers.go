package helpers

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mutx    sync.Mutex
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	name := logingClient(conn)
	fmt.Printf("[%s] has joined the server\n", name)

	// Here you can add logic to handle messages from the client
	for {
		// read the  message from  the client
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading from client: %v\n", err)
			return
		}
		// handle the message
		msg := string(buf)
		destributeMessages(conn, fmt.Sprintf("Message from [%s] : %s", name, msg))
	}
}

func logingClient(conn net.Conn) string {
	_, err := conn.Write([]byte("Enter your name: "))
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return ""
	}

	// Create a buffer for reading the name
	nameBuffer := make([]byte, 256)
	n, err := conn.Read(nameBuffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return ""
	}

	name := strings.TrimSpace(string(nameBuffer[:n]))
	name = checkName(conn, name)
	mutx.Lock()
	clients[conn] = name
	mutx.Unlock()
	destributeMessages(conn, fmt.Sprintf("[%s] has joined the chat.\n", name))
	return name
}

func checkName(conn net.Conn, name string) string {
	mutx.Lock()
	defer mutx.Unlock()

	for name == "" || isNameTaken(name) {
		_, err := conn.Write([]byte("The name is already taken or invalid, please enter a new name: "))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return ""
		}
		nameBuffer := make([]byte, 256)
		n, err := conn.Read(nameBuffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return ""
		}
		name = strings.TrimSpace(string(nameBuffer[:n]))
	}
	return name
}

func isNameTaken(name string) bool {
	for _, existingName := range clients {
		if existingName == name {
			return true
		}
	}
	return false
}

func destributeMessages(conn net.Conn, msg string) {
	mutx.Lock()
	defer mutx.Unlock()
	for conns := range clients {
		if conns != conn {
			_, err := conns.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error writing to client:", err)
				continue
			}
		}
	}
}
