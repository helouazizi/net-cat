package helpers

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mutx    sync.Mutex
)

/*
this function about to promot  the client in two parts
[first part] :  get the client name and chaecking for validation
already exist in the map or invalid one
after that  it will send the name to the server abd store it into the map
[second part] : promot  the client to enter the chat room and handle his  messages
*/

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// the first part began here
	name, err := logingClient(conn)
	if err != nil {
		fmt.Println("A client disconected befor providing a name")
		return
	} else {
		mutx.Lock()
		clients[conn] = name
		mutx.Unlock()
		conn.Write([]byte("Bienvenido! " + name + "\n"))
		fmt.Printf("[%s] has joined the server\n", name)
		fmt.Printf("[%s] added to the client list\n", name)
		destributeMessages(conn, fmt.Sprintf("[%s] has joined the chat.\n", name))
	}

	// the second part  began here
	// lets  start to handle the messages
	for {
		// read the  message from  the client
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("[%s] Disconected from the server\n", name)
			destributeMessages(conn, fmt.Sprintf("[%s]:leaved the chat\n", name))
			mutx.Lock()
			delete(clients, conn)
			mutx.Unlock()
			fmt.Printf("[%s] removed from the clients list\n", name)
			return
		}
		msg := string(buf)
		destributeMessages(conn, fmt.Sprintf("[%s] : %s", name, msg))
	}
}

/*
this function promot the user to  enter his name
and check if the name is valid or already exist
*/
func logingClient(conn net.Conn) (string, error) {
	Welcommessage := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n"

	_, err := conn.Write([]byte(Welcommessage + "[Enter your name]: "))
	if err != nil {
		return "", err
	}

	// Create a buffer for reading the name
	nameBuffer := make([]byte, 1024)
	n, err := conn.Read(nameBuffer)
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(string(nameBuffer[:n]))
	name, err = checkName(conn, name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func checkName(conn net.Conn, name string) (string, error) {
	mutx.Lock()
	defer mutx.Unlock()
	try := 1
	for name == "" || isNameTaken(name) {
		_, err := conn.Write([]byte("Invalid or taken name. Please choose another name: "))
		if err != nil {
			return "", err
		}
		nameBuffer := make([]byte, 1024)
		n, err := conn.Read(nameBuffer)
		if err != nil {
			return "", err
		}
		try++
		if try > 3 {
			_, err := conn.Write([]byte("You have tried 3 times. plaese try again later"))
			if err != nil {
				return "", err
			}
			conn.Close()
			os.Exit(1)
		}
		// why this n is fucking of the logic
		name = strings.TrimSpace(string(nameBuffer[:n]))
	}
	return name, nil
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
				continue
			}
		}
	}
}
