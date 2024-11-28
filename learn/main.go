package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	// Create a TCP socket
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	fmt.Println(sockfd)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		os.Exit(1)
	}
	defer syscall.Close(sockfd)

	// Bind the socket to an address
	addr := &syscall.SockaddrInet4{Port: 8080}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	if err := syscall.Bind(sockfd, addr); err != nil {
		fmt.Println("Error binding socket:", err)
		os.Exit(1)
	}

	// Listen for incoming connections
	if err := syscall.Listen(sockfd, 500); err != nil {
		fmt.Println("Error listening on socket:", err)
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 8080...")

	for {
		// Accept an incoming connection
		connfd, _, err := syscall.Accept(sockfd)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected")

		// Handle the connection
		go handleConnection(connfd)
	}
}

func handleConnection(connfd int) {
	defer syscall.Close(connfd)

	// Read data from the client
	buffer := make([]byte, 1024)
	n, err := syscall.Read(connfd, buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	fmt.Println("Received:", string(buffer[:n]))

	// Send a response back to the client
	_, err = syscall.Write(connfd, []byte("Hello from server!"))
	if err != nil {
		fmt.Println("Error writing to client:", err)
	}
}
