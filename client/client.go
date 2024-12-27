package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//fmt.Print("Enter server address (e.g., localhost:8080): ")
	fmt.Println("default address localhost 8080")
	var serverAddr string = "localhost:8080"
	//fmt.Scanln(&serverAddr)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start reading server messages in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Send messages to the server
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("your name for this chat >> ")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	conn.Write([]byte(msg + "\n"))

	for {
		fmt.Print(">> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg == "/quit" {
			fmt.Println("Exiting chat...")
			break
		}

		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
