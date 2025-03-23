package main

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rah-0/nabu"

	. "github.com/rah-0/hyperion/hconn"
	. "github.com/rah-0/hyperion/model"
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

type EntityConfig struct {
	Name string
}

type EntityStorage struct {
	Disk   *Disk
	Memory *register.Entity
}

type Node struct {
	Host            Host
	Path            Path
	errCh           chan error
	Status          NodeStatus
	HConn           *HConn
	Peers           []*Node
	Entities        []*EntityConfig
	EntitiesStorage []*EntityStorage

	Mu sync.Mutex
}

func NewNode(h Host, p Path, ecs []*EntityConfig) *Node {
	return &Node{
		Host:     h,
		Path:     p,
		Entities: ecs,
		errCh:    make(chan error, 1),
		Peers:    []*Node{},
	}
}

func ConnectToNodeWithHostAndPort(ip string, port string) (*HConn, error) {
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", ip+":"+port)
		if err == nil {
			return NewHConn(conn), nil
		}
		if strings.Contains(err.Error(), "connection refused") {
			nabu.FromMessage("trying to connect to: [" + ip + ":" + port + "]").Log()
		} else {
			return nil, err
		}
		time.Sleep(1 * time.Second)
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

func (x *Node) Start() error {
	if err := x.checkDataDir(); err != nil {
		return err
	}

	// Config per node targets an entity by name but here we find all versions for that entity
	for _, e := range x.Entities {
		for _, re := range register.Entities {
			if e.Name == re.Name {
				disk := NewDisk().WithPath(filepath.Join(x.Path.Data, re.DbFileName)).WithEntity(re)
				if err := disk.OpenFile(); err != nil {
					return err
				}

				x.EntitiesStorage = append(x.EntitiesStorage, &EntityStorage{
					Disk:   disk,
					Memory: re,
				})
			}
		}
	}

	if err := x.loadEntitiesFromDisk(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", x.getListenAddress())
	if err != nil {
		return nabu.FromError(err).WithArgs(x.Host).Log()
	}

	x.handleErrors()
	defer func() {
		if err = listener.Close(); err != nil {
			x.errCh <- err
		}
	}()
	go x.connectToPeers()
	x.acceptConnections(listener)

	return nil
}

func (x *Node) checkDataDir() error {
	exists, err := PathExists(x.Path.Data)
	if err != nil {
		return err
	}
	if !exists {
		err = DirectoryCreate(x.Path.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Node) loadEntitiesFromDisk() error {
	for _, s := range x.EntitiesStorage {
		d := s.Disk
		exists, err := PathExists(d.Path)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}

		if err = d.DataCleanup(); err != nil {
			return err
		}

		entities, err := d.DataReadAll()
		if err != nil {
			return err
		}

		s.Memory.New().MemorySet(entities)
	}

	return nil
}

func (x *Node) connectToPeers() {
	x.WaitStatusActive()

	var newPeers []*Node
	for _, node := range GlobalConfig.Nodes {
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
	return fmt.Sprintf("%s:%d", x.Host.Name, x.Host.Port)
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
		if err := hc.C.Close(); err != nil {
			nabu.FromError(err).Log()
		}
	}()

	nabu.FromMessage("new connection from [" + hc.C.RemoteAddr().String() + "] to [" + hc.C.LocalAddr().String() + "]").Log()

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
			err = hc.Send(msg)
		} else if msg.Type == MessageTypeInsert {
			for _, e := range x.EntitiesStorage {
				if e.Memory.Name == msg.Entity.Name && e.Memory.Version == msg.Entity.Version {
					entity := e.Memory.New()
					entity.SetBufferData(msg.Entity.Data)
					if err := entity.Decode(); err != nil {
						nabu.FromError(err).Log()
					}
					entity.MemoryAdd()

					if err := e.Disk.DataWrite(msg.Entity.Data); err != nil {
						nabu.FromError(err).Log()
					}

					if err := hc.Send(Message{
						String: "Sample Here",
					}); err != nil {
						nabu.FromError(err).Log()
					}
					break
				}
			}
		}

		if err != nil {
			nabu.FromError(err).Log()
			break
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
					nabu.FromError(err).WithLevelFatal().Log()
				}
			}
		}
	}(x)
}

func (x *Node) WaitStatusActive() {
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
