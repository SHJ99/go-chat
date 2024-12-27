package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	sender  string
	content string
}

func handleConnection(conn net.Conn, messages chan<- Message, disconnect chan<- string) {
	defer conn.Close()

	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)
	conn.Write([]byte("welcome to go-chatting room, " + name + "\n"))

	messages <- Message{sender: "Server", content: fmt.Sprintf("%s has joined the chat", name)}

	for { // Listen message
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break // Client disconnect
		}
		msg = strings.TrimSpace(msg)

		if msg == "/quit" {
			break
		}

		messages <- Message{sender: name, content: msg}
	}

	// Notify disconnection
	disconnect <- name
	messages <- Message{sender: "Server", content: fmt.Sprintf("%s has left the chat", name)}
}

func broadcastMessages(messages <-chan Message, clients map[string]net.Conn, newClient <-chan net.Conn, disconnect <-chan string) {
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				conn.Write([]byte(fmt.Sprintf("%s: %s\n", msg.sender, msg.content)))
			}
		case conn := <-newClient:
			clients[conn.RemoteAddr().String()] = conn
		case name := <-disconnect:
			delete(clients, name)
		}
	}
}

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port", port)

	clients := make(map[string]net.Conn)
	messages := make(chan Message)
	newClient := make(chan net.Conn)
	disconnect := make(chan string)

	go broadcastMessages(messages, clients, newClient, disconnect)

	for { // Accept new clients
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		newClient <- conn
		go handleConnection(conn, messages, disconnect)
	}
}
