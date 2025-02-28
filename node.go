package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/rah-0/nabu"
)

type NodeStatus int

const (
	NodeStatusStarting NodeStatus = iota
	NodeStatusReady
	NodeStatusReadyToReceiveData
)

var wg sync.WaitGroup

func startNodes() {
	errCh := make(chan error, 1)

	for _, node := range nodes {
		node.errCh = errCh
		node.checkDataDir()

		wg.Add(1)
		go node.listenPortForStatus()
	}

	waitNodesStatusPortToBeReady()
	go func() {
		wg.Wait()
		close(errCh)
	}()
}

func waitNodesStatusPortToBeReady() {
	for {
		allReady := true
		for _, node := range nodes {
			status := node.Status
			if status != NodeStatusReady {
				allReady = false
				break
			}
		}
		if allReady {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (x *Node) checkDataDir() {
	exists, err := pathExists(x.Path.Data)
	if err != nil {
		x.errCh <- err
		return
	}
	if !exists {
		err := pathCreateDirs(x.Path.Data)
		if err != nil {
			x.errCh <- err
			return
		}
	}
}

func (x *Node) listenPortForStatus() {
	defer wg.Done()

	port := x.Host.Ports.Status
	if portIsInUse(port) {
		x.errCh <- nabu.FromError(ErrConfigNodePortStatusNotAvailable).WithArgs(port).Log()
		return
	}

	listener, err := net.Listen("tcp", x.getListenAddressForStatus())
	if err != nil {
		x.errCh <- err
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			x.errCh <- err
		}
	}()

	x.mu.Lock()
	x.Status = NodeStatusReady
	x.mu.Unlock()
	for {
		conn, err := listener.Accept()
		if err != nil {
			x.errCh <- err
			return
		}

		go x.handleConnection(conn)
	}
}

func (x *Node) getListenAddressForStatus() string {
	return fmt.Sprintf(":%d", x.Host.Ports.Status)
}

func (x *Node) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			nabu.FromError(err).Log()
		}
	}()

	fmt.Println("New client connected:", conn.RemoteAddr())

	// Read data from client
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	fmt.Println("Received:", string(buf[:n]))

	// Example: Send a response
	conn.Write([]byte("Hello from server!\n"))
}
