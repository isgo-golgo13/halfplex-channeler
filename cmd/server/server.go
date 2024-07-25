package main

import (
	"fmt"
	"net"
	"os"

	"github.com/isgo-golgo13/fifochannel/svckit"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:8000" // Default to 0.0.0.0:8000 if not set
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fifo := &svckit.FifoSocketChannel{Conn: conn}

	// Receiving data
	recvBuf := make([]byte, 64)
	if err, n := fifo.Recv(recvBuf, int64(len(recvBuf))); err != nil {
		fmt.Println("Recv error:", err)
		return
	} else {
		fmt.Println("Received bytes:", n)
		fmt.Println("Received data:", string(recvBuf[:n]))
	}

	// Sending response
	sendBuf := []byte("hello from server")
	if err, n := fifo.Send(sendBuf, int64(len(sendBuf))); err != nil {
		fmt.Println("Send error:", err)
	} else {
		fmt.Println("Sent bytes:", n)
	}
}
