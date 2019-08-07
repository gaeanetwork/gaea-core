package server

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	server := NewTeeServer(":12345")
	go server.Start()
	for {
		conn, err := net.Dial("tcp", server.address)
		if err == nil {
			conn.Close()
			break
		}
	}
	assert.NotNil(t, server.Server())

	// Repeated start
	err := server.Start()
	assert.Error(t, err, "bind: address already in use")

	// Stop
	server.GracefulStop()
	_, err = net.Dial("tcp", server.address)
	assert.Error(t, err, "connect: connection refused")

	// Invalid address - missing port
	server = NewTeeServer("127.0.0.1")
	err = server.Start()
	assert.Error(t, err, "missing port in address")

	// Invalid address - invalid port
	server = NewTeeServer("127.0.0.1:12345789")
	err = server.Start()
	assert.Error(t, err, "invalid port")

	// Invalid address - no such host
	server = NewTeeServer("127.0.0.256:12345")
	err = server.Start()
	assert.Error(t, err, "no such host")
}
