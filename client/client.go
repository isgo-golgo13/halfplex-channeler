package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"halfplex-channeler/svckit"
)

func main() {
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	ctx := context.Background()
	channel := svckit.NewHalfPlexChanneler(conn, conn)

	// Sending a message to the server
	data := strings.NewReader("Hello from Client!")
	if err, _ := channel.Send(ctx, data, int64(data.Len())); err != nil {
		log.Fatal("Send error:", err)
	}

	// Receiving a response from the server
	var responseBuffer bytes.Buffer
	if err, _ := channel.Recv(ctx, &responseBuffer, 1024); err != nil {
		log.Fatal("Receive error:", err)
	}
	fmt.Println("Received from server:", responseBuffer.String())

	// Sending a timed message
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := channel.SendTimeout(timeoutCtx, strings.NewReader("Timed Message from Client"), 1024, 5*time.Second); err != nil {
		log.Fatal("SendTimeout error:", err)
	}
}
