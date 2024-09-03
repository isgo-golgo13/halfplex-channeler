package svckit

import (
	"context"
	"io"
	"net"
	"sync"
	"time"
)

// HalfPlexChannel defines the interface for half-duplex communication with context support.
type HalfPlexChannel interface {
	Send(ctx context.Context, data io.Reader, n int64) (int64, error)
	SendTimeout(ctx context.Context, data io.Reader, n int64, timeout time.Duration) (int64, error)
	Recv(ctx context.Context, writer io.Writer, n int64) (int64, error)
	SendAll(ctx context.Context, data []io.Reader, n int64) (int64, error)
	Close(ctx context.Context) error
}

// HalfPlexChanneler implements the HalfPlexChannel interface.
type HalfPlexChanneler struct {
	conn   net.Conn
	closed bool
	mu     sync.Mutex // Mutex to enforce half-duplex behavior
}

// NewHalfPlexChanneler creates a new HalfPlexChanneler.
func NewHalfPlexChanneler(conn net.Conn) *HalfPlexChanneler {
	return &HalfPlexChanneler{conn: conn}
}

// Send sends data from an io.Reader to the connection.
func (hc *HalfPlexChanneler) Send(ctx context.Context, data io.Reader, n int64) (int64, error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.closed {
		return 0, io.ErrClosedPipe
	}

	sentBytes, err := io.CopyN(hc.conn, data, n)
	return sentBytes, err
}

// SendTimeout sends data with a timeout.
func (hc *HalfPlexChanneler) SendTimeout(ctx context.Context, data io.Reader, n int64, timeout time.Duration) (int64, error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.closed {
		return 0, io.ErrClosedPipe
	}

	done := make(chan struct{})
	var sentBytes int64
	var sendErr error

	go func() {
		sentBytes, sendErr = hc.Send(ctx, data, n)
		close(done)
	}()

	select {
	case <-done:
		return sentBytes, sendErr
	case <-ctx.Done():
		return sentBytes, ctx.Err()
	case <-time.After(timeout):
		return sentBytes, context.DeadlineExceeded
	}
}

// Recv receives data into an io.Writer from the connection.
func (hc *HalfPlexChanneler) Recv(ctx context.Context, writer io.Writer, n int64) (int64, error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.closed {
		return 0, io.ErrClosedPipe
	}

	recvBytes, err := io.CopyN(writer, hc.conn, n)
	return recvBytes, err
}

// SendAll sends all data from a slice of io.Reader objects.
func (hc *HalfPlexChanneler) SendAll(ctx context.Context, data []io.Reader, n int64) (int64, error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.closed {
		return 0, io.ErrClosedPipe
	}

	var totalSent int64
	for _, reader := range data {
		sent, err := hc.Send(ctx, reader, n-totalSent)
		totalSent += sent
		if err != nil {
			return totalSent, err
		}
		if totalSent >= n {
			break
		}
	}
	return totalSent, nil
}

// Close closes the connection.
func (hc *HalfPlexChanneler) Close(ctx context.Context) error {
	hc.closed = true
	return hc.conn.Close()
}
