package main

import (
	"fmt"
	"net"
	"os"

	"netcat/helpers"
)

func main() {
	// in first lets promt the ussr to provide a valdi port
	// or using our default port 8989
	var (
		port string
		err  error
	)

	if len(os.Args) < 2 {
		port = "8989"
	} else {
		port = os.Args[1]
		port, err = helpers.CheckPrort(port)
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
	}
	// lets craete a tcp server using  net package
	// we will use net.Listen() function to listen on the default port 8989 if there is no  port specified
	listner, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error  starting on port" + port)
		return
	}
	defer listner.Close()
	fmt.Printf("Server is listening on port %s...", port)

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
