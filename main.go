package main

import (
	"fmt"
	"net"

	"netcat/helpers"
)

func main() {
	// lets craete a tcp server using  net package
	// we will use net.Listen() function to listen on the default port 8989 if there is no  port specified
	listner, err := net.Listen("tcp", ":8989")
	if err != nil {
		fmt.Println("Error  starting on port 8989")
		return
	}
	defer listner.Close()
	fmt.Println("Server is listening on port 8989...")

	//  we will use a for loop to accept incoming connections
	for {
		// we will use net.conn.Accept() function to accept incoming connections
		// because we use tcp  protocol, that is an oriented connection protocol
		// we need to establish  a connection before we can send data or  receive data
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
			return
		}

		// know  that we have a connection, we can start reading from it or  writing to it
		// but we have a broblem how to handle the connections concurently
		// we can use goroutines to handle the connections concurrently
		go helpers.HandleClient(conn)
	}
}
