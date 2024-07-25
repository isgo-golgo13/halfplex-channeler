package main

import (
	"fmt"
	"net"
	"os"

	"github.com/isgo-golgo13/fifochannel/svckit"
)

func main() {
	serverAddr := "server:8000"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fifo := &svckit.FifoSocketChannel{Conn: conn}

	// Sending data
	sendBuf := []byte("hello from client")
	if err, n := fifo.Send(sendBuf, int64(len(sendBuf))); err != nil {
		fmt.Println("Send error:", err)
	} else {
		fmt.Println("Sent bytes:", n)
	}

	// Receiving data
	recvBuf := make([]byte, 64)
	if err, n := fifo.Recv(recvBuf, int64(len(recvBuf))); err != nil {
		fmt.Println("Recv error:", err)
	} else {
		fmt.Println("Received bytes:", n)
		fmt.Println("Received data:", string(recvBuf[:n]))
	}
}
