package util

import (
	"net"
	"strconv"
)

// GetAvailablePort gets an available port by asking the OS for a free port
// by binding to port 0 and then closing the connection, returning the port that was assigned
func GetAvailablePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}
	
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer l.Close()
	
	return l.Addr().(*net.TCPAddr).Port
}

// GetAvailablePortStr returns an available port as a string
func GetAvailablePortStr() string {
	return strconv.Itoa(GetAvailablePort())
}
