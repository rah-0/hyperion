package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/rah-0/nabu"
)

type NodeStatus int

const (
	NodeStatusStarting NodeStatus = iota
	NodeStatusReady
	NodeStatusReadyToReceiveData
)

var wg sync.WaitGroup

type Node struct {
	Host   Host
	Path   Path
	errCh  chan error
	Status NodeStatus

	mu sync.Mutex
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

	port := x.Host.Port
	if portIsInUse(port) {
		x.errCh <- nabu.FromError(ErrConfigNodePortStatusNotAvailable).WithArgs(port).Log()
		return
	}

	listener, err := net.Listen("tcp", x.getListenAddress())
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

		c := NewHConn(conn)
		go x.handleConnection(c)
	}
}

func (x *Node) getListenAddress() string {
	return fmt.Sprintf(":%d", x.Host.Port)
}

func (x *Node) handleConnection(hc *HConn) {
	defer func() {
		if err := hc.c.Close(); err != nil {
			nabu.FromError(err).Log()
		}
	}()

	for {
		msg, err := hc.Receive()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			nabu.FromError(err).Log()
			break
		}

		if msg.Type == MessageTypeTest {
			msg.String += "Received"
		}

		if msg.Mode == ModeSync {
			err = hc.Send(msg)
			if err != nil {
				nabu.FromError(err).Log()
				break
			}
		}
	}
}
