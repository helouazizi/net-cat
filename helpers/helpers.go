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
	name, err := logingClient(conn)
	if err != nil {
		fmt.Println("A client disconected befor providing a name")
		return
	} else {
		mutx.Lock()
		clients[conn] = name
		mutx.Unlock()

		fmt.Printf("[%s] has joined the server\n", name)
		fmt.Printf("[%s] added to the client list\n", name)
		destributeMessages(conn, fmt.Sprintf("[%s] has joined the chat.\n", name))
	}

	// Lets  handle messages from the client
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

func logingClient(conn net.Conn) (string, error) {
	// lets send a message to the client to provide a name
	Welcommessage := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n"

	_, err := conn.Write([]byte(Welcommessage + "[Enter your name]: "))
	if err != nil {
		return "", err
	}

	// Create a buffer for reading the name
	nameBuffer := make([]byte, 256)
	n, err := conn.Read(nameBuffer)
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(string(nameBuffer[:n]))
	name = checkName(conn, name)
	return name, nil
}

func checkName(conn net.Conn, name string) string {
	mutx.Lock()
	defer mutx.Unlock()

	for name == "" || isNameTaken(name) {
		_, err := conn.Write([]byte("Invalid or taken name. Please choose another name: "))
		if err != nil {
			return ""
		}
		nameBuffer := make([]byte, 256)
		n, err := conn.Read(nameBuffer)
		if err != nil {
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
				continue
			}
		}
	}
}
