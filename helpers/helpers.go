package helpers

import (
	"bufio"
	"fmt"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// in first lets promt the client to enter his name
	_, err := writer.WriteString("Enter your  name: ")
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}
	//  flush the buffer to ensure  the message is sent imediatly
	writer.Flush()
	name, _ := reader.ReadString('\n')
	// lets remove the  \n from the name
	name = name[:len(name)-1]
	fmt.Printf("[%s]joined the server", name)
	fmt.Println()

	// ok know  lets promt the client to enter his messages
	for {

		_, err := writer.WriteString("Enter your  message: ")
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
		writer.Flush()
		
		// lets read the message from the client with \n as  the delimiter
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error  reading from client:", err)
			return
		}
		fmt.Printf("Received message from  %s: %s", name, msg)

		//  Send response back to client
		_, err = writer.WriteString("Server response : hello from  server\n")
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
		writer.Flush()
	}
}
