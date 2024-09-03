package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"halfplex-channeler/svckit"
)

// Server struct represents a TCP server with a configurable port.
type Server struct {
	port string
}

// Option defines a functional option for configuring the server.
type Option func(*Server)

// WithPort is a functional option to set the port for the server.
func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// NewServer creates a new Server with the given options.
func NewServer(opts ...Option) *Server {
	s := &Server{port: "8080"}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {
	// Initialize the server with a configurable port using functional options.
	server := NewServer(WithPort("8080"))

	// Listen for incoming connections on the specified port.
	ln, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer ln.Close()

	// Channel to handle OS signals for graceful shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Context for managing connection handling and graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown.
	go func() {
		<-sigs
		cancel()
		ln.Close()
		fmt.Println("Shutting down server...")
	}()

	// Accept incoming connections in a loop.
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return // Exit loop on server shutdown.
			default:
				log.Println("Failed to accept connection:", err)
			}
			continue
		}

		// Handle each connection in a separate goroutine.
		go handleConnection(ctx, conn)
	}
}

// handleConnection handles the communication with a connected client.
func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Initialize a HalfPlexChanneler for communication.
	channel := svckit.NewHalfPlexChanneler(conn, conn)

	// Receive a message from the client.
	var requestBuffer bytes.Buffer
	if err, _ := channel.Recv(ctx, &requestBuffer, 1024); err != nil {
		log.Println("Receive error:", err)
		return
	}
	fmt.Println("Received from client:", requestBuffer.String())

	// Send a response back to the client.
	response := strings.NewReader("Hello from Server!")
	if err, _ := channel.Send(ctx, response, int64(response.Len())); err != nil {
		log.Println("Send error:", err)
	}
}
