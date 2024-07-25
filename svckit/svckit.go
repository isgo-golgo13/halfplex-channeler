package svckit

import (
	"fmt"
	"net"
	"time"
)

// FifoChannel interface with the specified methods
type FifoChannel interface {
	Send(buf []byte, n int64) (error, int64)
	SendTimeout(buf []byte, n int64, timeout <-chan time.Duration) (error, int64)
	Recv(buf []byte, n int64) (error, int64)
}

// FifoSocketChannel struct that will implement the FifoChannel interface
type FifoSocketChannel struct {
	Conn net.Conn
}

// Implement the Send method for FifoSocketChannel
func (f *FifoSocketChannel) Send(buf []byte, n int64) (error, int64) {
	_, err := f.Conn.Write(buf[:n])
	if err != nil {
		return err, 0
	}
	return nil, n
}

// Implement the SendTimeout method for FifoSocketChannel
func (f *FifoSocketChannel) SendTimeout(buf []byte, n int64, timeout <-chan time.Duration) (error, int64) {
	select {
	case <-timeout:
		return fmt.Errorf("send timeout"), 0
	default:
		return f.Send(buf, n)
	}
}

// Implement the Recv method for FifoSocketChannel
func (f *FifoSocketChannel) Recv(buf []byte, n int64) (error, int64) {
	readBytes, err := f.Conn.Read(buf[:n])
	if err != nil {
		return err, 0
	}
	return nil, int64(readBytes)
}
