package svckit

import (
	"context"
	"io"
	"time"
)

// HalfPlexChannel defines an interface for half-duplex communication with cancellable operations using context.
type HalfPlexChannel interface {
	Send(ctx context.Context, data io.Reader, n int64) (error, int64)
	SendTimeout(ctx context.Context, data io.Reader, n int64, timeout <-chan time.Duration) (error, int64)
	Recv(ctx context.Context, writer io.Writer, n int64) (error, int64)
	SendAll(ctx context.Context, data []io.Reader, n int64) (error, int64)
	Close(ctx context.Context) error
}

// HalfPlexChanneler is a concrete implementation of HalfPlexChannel.
type HalfPlexChanneler struct {
	sendWriter io.Writer
	recvReader io.Reader
	closed     bool
}

// NewHalfPlexChanneler creates a new HalfPlexChanneler.
func NewHalfPlexChanneler(sendWriter io.Writer, recvReader io.Reader) *HalfPlexChanneler {
	return &HalfPlexChanneler{
		sendWriter: sendWriter,
		recvReader: recvReader,
		closed:     false,
	}
}

// Send sends data with the provided context.
func (hc *HalfPlexChanneler) Send(ctx context.Context, data io.Reader, n int64) (error, int64) {
	if hc.closed {
		return io.ErrClosedPipe, 0
	}

	written := int64(0)
	buffer := make([]byte, 4096)

	for written < n {
		select {
		case <-ctx.Done():
			return ctx.Err(), written
		default:
			toRead := min(4096, n-written)
			bytesRead, err := data.Read(buffer[:toRead])
			if err != nil && err != io.EOF {
				return err, written
			}
			bytesWritten, err := hc.sendWriter.Write(buffer[:bytesRead])
			written += int64(bytesWritten)
			if err != nil {
				return err, written
			}
			if bytesRead == 0 {
				break
			}
		}
	}

	return nil, written
}

// SendTimeout sends data with a timeout using the provided context.
func (hc *HalfPlexChanneler) SendTimeout(ctx context.Context, data io.Reader, n int64, timeout time.Duration) (int64, error) {
	if hc.closed {
		return 0, io.ErrClosedPipe
	}

	done := make(chan struct{})
	var sentBytes int64
	var sendErr error

	go func() {
		sendErr, sentBytes = hc.Send(ctx, data, n)
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

// Recv receives data with the provided context.
func (hc *HalfPlexChanneler) Recv(ctx context.Context, writer io.Writer, n int64) (error, int64) {
	if hc.closed {
		return io.ErrClosedPipe, 0
	}

	read := int64(0)
	buffer := make([]byte, 4096)

	for read < n {
		select {
		case <-ctx.Done():
			return ctx.Err(), read
		default:
			toWrite := min(4096, n-read)
			bytesRead, err := hc.recvReader.Read(buffer[:toWrite])
			if err != nil && err != io.EOF {
				return err, read
			}
			bytesWritten, err := writer.Write(buffer[:bytesRead])
			read += int64(bytesWritten)
			if err != nil {
				return err, read
			}
			if bytesRead == 0 {
				break
			}
		}
	}

	return nil, read
}

// SendAll sends data from multiple readers using the provided context.
func (hc *HalfPlexChanneler) SendAll(ctx context.Context, data []io.Reader, n int64) (error, int64) {
	if hc.closed {
		return io.ErrClosedPipe, 0
	}

	var totalSent int64

	for _, reader := range data {
		err, sent := hc.Send(ctx, reader, n-totalSent)
		totalSent += sent
		if err != nil {
			return err, totalSent
		}
		if totalSent >= n {
			break
		}
	}

	return nil, totalSent
}

// Close closes the HalfPlexChanneler using the provided context.
func (hc *HalfPlexChanneler) Close(ctx context.Context) error {
	hc.closed = true
	// Implement additional closing logic if needed
	return nil
}

// min returns the smaller of x or y.
func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
