package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/isgo-golgo13/fifochannel/svckit"
	"github.com/joho/godotenv"
)

const CONFIG_FILE = "../../.env"

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(CONFIG_FILE)
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

	// Create a channel to listen for interrupt or terminate signals from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal that the server is done
	doneChan := make(chan bool, 1)

	// Handle incoming connections
	go func() {
		fmt.Println("Server is listening on", listenAddr)
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-doneChan:
					return
				default:
					fmt.Println("Failed to accept connection:", err)
				}
				continue
			}

			go handleConnection(conn)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nShutting down server...")

	// Signal that the server should stop accepting new connections
	close(doneChan)

	// Close the listener
	listener.Close()
	fmt.Println("Server stopped gracefully.")
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
	sendBuf := []byte("ACK of client sent packket from server")
	if err, n := fifo.Send(sendBuf, int64(len(sendBuf))); err != nil {
		fmt.Println("Send error:", err)
	} else {
		fmt.Println("Sent bytes:", n)
	}
}
