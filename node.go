package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rah-0/nabu"

	"github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

type NodeStatus int

const (
	NodeStatusStarting NodeStatus = iota
	NodeStatusActive
	NodeStatusReady
)

type Path struct {
	Data string // Where data will be stored
}

type Host struct {
	Name string
	Port int
}

type Node struct {
	Host     Host
	Path     Path
	errCh    chan error
	Status   NodeStatus
	HConn    *HConn
	Peers    []*Node
	Entities []*register.Entity

	Mu sync.Mutex
}

func NewNode(h Host, p Path, es []*register.Entity) *Node {
	return &Node{
		Host:     h,
		Path:     p,
		Entities: es,
		errCh:    make(chan error, 1),
		Peers:    []*Node{},
	}
}

func ConnectToNode(x *Node) (*HConn, error) {
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", x.getListenAddress())
		if err == nil {
			return NewHConn(conn), nil
		}
		if strings.Contains(err.Error(), "connection refused") {
			nabu.FromMessage("trying to connect to: [" + x.getListenAddress() + "]").Log()
		} else {
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}
}

func (x *Node) checkDataDir() {
	exists, err := PathExists(x.Path.Data)
	if err != nil {
		x.errCh <- err
		return
	}
	if !exists {
		err := DirectoryCreate(x.Path.Data)
		if err != nil {
			x.errCh <- err
			return
		}
	}
}

func (x *Node) Start() {
	x.handleErrors()
	x.checkDataDir()

	listener, err := net.Listen("tcp", x.getListenAddress())
	if err != nil {
		x.errCh <- nabu.FromError(err).WithArgs(x.Host)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			x.errCh <- err
		}
	}()

	go x.connectToPeers()
	x.acceptConnections(listener)
}

func (x *Node) connectToPeers() {
	x.waitStatusActive()

	var newPeers []*Node
	for _, node := range config.Nodes {
		// Skip self
		if node.Host.Name == x.Host.Name && node.Host.Port == x.Host.Port {
			continue
		}

		c, err := ConnectToNode(node)
		if err != nil {
			x.errCh <- nabu.FromError(err).Log()
			continue
		}

		node.Mu.Lock()
		node.HConn = c
		node.Mu.Unlock()
		newPeers = append(newPeers, node)
	}

	// Update the peer list
	x.Mu.Lock()
	x.Peers = newPeers
	x.Status = NodeStatusReady
	x.Mu.Unlock()
}

func (x *Node) getListenAddress() string {
	return fmt.Sprintf(":%d", x.Host.Port)
}

func (x *Node) acceptConnections(listener net.Listener) {
	x.Mu.Lock()
	x.Status = NodeStatusActive
	x.Mu.Unlock()

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

func (x *Node) handleConnection(hc *HConn) {
	defer func() {
		if err := hc.c.Close(); err != nil {
			nabu.FromError(err).Log()
		}
	}()

	nabu.FromMessage("new connection from [" + hc.c.RemoteAddr().String() + "] to [" + hc.c.LocalAddr().String() + "]").Log()

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

func (x *Node) handleErrors() {
	go func(n *Node) {
		for {
			select {
			case err, ok := <-n.errCh:
				if !ok {
					return
				}
				if err != nil {
					if strings.Contains(err.Error(), "address already in use") {
						nabu.FromError(err).WithLevelFatal().Log()
						os.Exit(1)
					} else {
						nabu.FromError(err).Log()
					}
				}
			}
		}
	}(x)
}

func (x *Node) waitStatusActive() {
	for {
		x.Mu.Lock()
		status := x.Status
		x.Mu.Unlock()
		if status == NodeStatusActive {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
