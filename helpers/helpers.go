package helpers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// in first lets prompt the client to enter his name
	name := logingClient(reader, writer) 
	fmt.Printf("[%s]joined the server", name)
	fmt.Println()

	// ok know  lets promt the client to enter his messages
	destributeMessages(reader, writer, name)
}





/*
this function  is used to get the name of the client 
that trying to  connect to the server and follow  the rules of the chat
by providing  the name of the client if not prompt him  to enter a name
obligtory using the checkName  function
*/
func logingClient(reader *bufio.Reader, writer *bufio.Writer) string {
	_, err := writer.WriteString("Enter your  name: ")
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return ""
	}
	//  flush the buffer to ensure  the message is sent imediatly
	writer.Flush()

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	name = checkName(reader, writer, name)
	return name
}





/*
this function  is used to check if the name of the client is valid
or empty   if not prompt him  until  he enter a valid name

*/
func checkName(reader *bufio.Reader, writer *bufio.Writer, name string) string {
	for name == "" {

		_, err := writer.WriteString("Please provide a name : ")
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return ""
		}
		writer.Flush()
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
	}
	return name
}


/*
This function about destrubuting  the messages 
To the  all active clients and manage the client  list
by addin and  removing the clients from the list
*/

func destributeMessages(reader *bufio.Reader, writer *bufio.Writer, name string) {
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
	} // end of for loop
}
