package node

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rah-0/nabu"

	"github.com/rah-0/hyperion/disk"
	"github.com/rah-0/hyperion/hconn"
	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/register"
	"github.com/rah-0/hyperion/template"
	"github.com/rah-0/hyperion/util"
)

type Status int

const (
	StatusStarting Status = iota
	StatusActive
	StatusReady
	StatusShutdown
)

type EntityStorage struct {
	Disk   *disk.Disk
	Memory *register.Entity
}

type Entity struct {
	Name string
}

type Path struct {
	Data string
}

type Host struct {
	Name string
	IP   string
	Port int
}

type Node struct {
	// Props coming from json config
	Host     Host
	Path     Path
	Entities []Entity

	ErrCh           chan error
	Status          Status
	HConn           *hconn.HConn
	Peers           []*Node
	EntitiesStorage []*EntityStorage
	PeerConnected   bool

	Mu sync.Mutex
}

func NewNode() *Node {
	return &Node{
		ErrCh: make(chan error, 1),
	}
}

func (x *Node) WithHost(name string, ip string, port int) *Node {
	x.Host.Name = name
	x.Host.IP = ip
	x.Host.Port = port
	return x
}

func (x *Node) WithPath(pathData string) *Node {
	x.Path.Data = pathData
	return x
}

func (x *Node) AddEntity(name string) *Node {
	x.Entities = append(x.Entities, Entity{Name: name})
	return x
}

func (x *Node) AddPeer(p *Node) *Node {
	x.Peers = append(x.Peers, p)
	return x
}

func ConnectToNodeWithHostAndPort(ip string, port string) (*hconn.HConn, error) {
	if err := template.RegisterEntities(); err != nil {
		return nil, err
	}

	for {
		conn, err := net.DialTimeout("tcp", ip+":"+port, 5*time.Second)
		if err == nil {
			return hconn.NewHConn(conn), nil
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "i/o timeout") {
			nabu.FromMessage("trying to connect to: [" + ip + ":" + port + "]").Log()
		} else {
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}
}

func ConnectToNode(x *Node) (*hconn.HConn, error) {
	if err := template.RegisterEntities(); err != nil {
		return nil, err
	}

	for {
		conn, err := net.DialTimeout("tcp", x.getListenAddress(), 5*time.Second)
		if err == nil {
			return hconn.NewHConn(conn), nil
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "i/o timeout") {
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
			if e.Name == re.EntityBase.Name {
				d := disk.NewDisk().WithPath(filepath.Join(x.Path.Data, re.EntityBase.DbFileName)).WithEntity(re)
				if err := d.OpenFile(); err != nil {
					return err
				}

				x.EntitiesStorage = append(x.EntitiesStorage, &EntityStorage{
					Disk:   d,
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
			x.ErrCh <- err
		}
	}()
	go x.connectToPeers()
	x.acceptConnections(listener)

	return nil
}

func (x *Node) checkDataDir() error {
	exists, err := util.PathExists(x.Path.Data)
	if err != nil {
		return err
	}
	if !exists {
		err = util.DirectoryCreate(x.Path.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Node) loadEntitiesFromDisk() error {
	for _, s := range x.EntitiesStorage {
		d := s.Disk
		exists, err := util.PathExists(d.Path)
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

		for _, e := range entities {
			e.MemoryAdd()
		}
	}

	return nil
}

func (x *Node) connectToPeers() {
	x.WaitStatusActive()

	for _, node := range x.Peers {
		c, err := ConnectToNode(node)
		if err != nil {
			x.ErrCh <- nabu.FromError(err).Log()
			continue
		}

		node.Mu.Lock()
		node.PeerConnected = true
		node.HConn = c
		node.Mu.Unlock()
	}

	x.Mu.Lock()
	x.Status = StatusReady
	x.Mu.Unlock()
}

func (x *Node) getListenAddress() string {
	return fmt.Sprintf("%s:%d", x.Host.IP, x.Host.Port)
}

func (x *Node) acceptConnections(listener net.Listener) {
	x.Mu.Lock()
	x.Status = StatusActive
	x.Mu.Unlock()

	for {
		conn, err := listener.Accept()
		if err != nil {
			x.ErrCh <- err
			return
		}

		c := hconn.NewHConn(conn)
		go x.handleConnection(c)
	}
}

func (x *Node) handleConnection(hc *hconn.HConn) {
	defer func() {
		if err := hc.Close(); err != nil {
			x.ErrCh <- nabu.FromError(err).Log()
		}
	}()

	nabu.FromMessage("new connection from [" + hc.C.RemoteAddr().String() + "] to [" + hc.C.LocalAddr().String() + "]").Log()

	for {
		msgIn, err := hc.Receive()
		if err != nil {
			x.ErrCh <- nabu.FromError(err).Log()
			break
		}

		msgOut := model.Message{}

		// Check node status before processing messages
		x.Mu.Lock()
		nodeStatus := x.Status
		x.Mu.Unlock()

		if nodeStatus == StatusShutdown {
			err := model.ErrNodeShutdown
			nabu.FromError(err).Log()
			msgOut.Error(err.Error())
			if err := hc.Send(msgOut); err != nil {
				x.ErrCh <- nabu.FromError(err).Log()
			}
			break
		}

		// Process message if node is not shutting down
		switch msgIn.Type {
		case model.MessageTypeInsert, model.MessageTypeDelete, model.MessageTypeUpdate:
			e := x.findEntityStorage(msgIn.Entity.Version, msgIn.Entity.Name)
			if e == nil {
				msgOut.Error("entity not found: [" + msgIn.Entity.Name + "]")
				break
			}

			entity := e.Memory.EntityExtension.New()
			entity.SetBufferData(msgIn.Entity.Data)
			if err := entity.Decode(); err != nil {
				msgOut.Error(err.Error())
				break
			}

			switch msgIn.Type {
			case model.MessageTypeInsert:
				entity.MemoryAdd()
			case model.MessageTypeDelete:
				entity.MemoryRemove()
			case model.MessageTypeUpdate:
				entity.MemoryUpdate()
			}

			if err = e.Disk.DataWrite(msgIn.Entity.Data); err != nil {
				msgOut.Error(err.Error())
				break
			}
			msgOut.Status = model.StatusSuccess

		case model.MessageTypeGetAll:
			e := x.findEntityStorage(msgIn.Entity.Version, msgIn.Entity.Name)
			if e == nil {
				msgOut.Error("entity not found: [" + msgIn.Entity.Name + "]")
				break
			}
			msgOut.Status = model.StatusSuccess
			msgOut.Models = e.Memory.EntityExtension.New().MemoryGetAll()

		case model.MessageTypeTest:
			msgOut.String = msgIn.String + "Received"

		case model.MessageTypeQuery:
			e := x.findEntityStorage(msgIn.Entity.Version, msgIn.Entity.Name)
			if e == nil {
				msgOut.Error("entity not found: [" + msgIn.Entity.Name + "]")
				break
			}

			r, err := e.HandleQuery(msgIn.Query)
			if err != nil {
				msgOut.Error(err.Error())
				break
			}

			msgOut.Status = model.StatusSuccess
			msgOut.Models = r
		}

		if err = hc.Send(msgOut); err != nil {
			x.ErrCh <- nabu.FromError(err).Log()
		}
	}
}

func (x *Node) findEntityStorage(version, name string) *EntityStorage {
	for _, e := range x.EntitiesStorage {
		if e.Memory.EntityBase.Version == version && e.Memory.EntityBase.Name == name {
			return e
		}
	}
	return nil
}

func (x *Node) handleErrors() {
	go func(n *Node) {
		for {
			n.Mu.Lock()
			// If channel is nil, it means the node is shutting down
			if n.ErrCh == nil {
				n.Mu.Unlock()
				return
			}

			// Create a local reference to the channel while holding the lock
			errCh := n.ErrCh
			n.Mu.Unlock()

			// Simple channel receive instead of select with single case
			err, ok := <-errCh
			if !ok {
				return
			}
			if err != nil {
				nabu.FromError(err).WithLevelFatal().Log()
			}
		}
	}(x)
}

func (x *Node) WaitStatusActive() {
	for {
		time.Sleep(100 * time.Millisecond)
		x.Mu.Lock()
		s := x.Status
		x.Mu.Unlock()
		if s == StatusActive {
			break
		}
	}
}

// Shutdown performs a graceful shutdown of the node, ensuring that all data is
// properly flushed to disk and resources are released.
// It closes all entity storage disks and cleans up resources
func (x *Node) Shutdown() error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	if x.Status == StatusShutdown {
		return nil
	}

	nabu.FromMessage("Shutting down node").WithArgs(x.Host).Log()

	x.Status = StatusShutdown
	if len(x.Peers) > 0 {
		x.Peers = nil
	}
	if x.HConn != nil {
		if err := x.HConn.Close(); err != nil {
			nabu.FromError(err).Log()
		}
		x.HConn = nil
	}
	for _, es := range x.EntitiesStorage {
		if es.Disk != nil {
			nabu.FromMessage("Closing entity storage disk").WithArgs(es.Disk.Path).Log()
			if err := es.Disk.Close(); err != nil {
				nabu.FromError(err).Log()
			}
			es.Disk = nil
		}
	}

	// Force cleanup any other references
	x.EntitiesStorage = nil

	// Close error channel - mutex is already locked in Shutdown
	if x.ErrCh != nil {
		errCh := x.ErrCh
		x.ErrCh = nil // Set to nil first to signal handleErrors goroutine to stop

		select {
		case <-errCh: // Channel already closed
			// Do nothing
		default:
			close(errCh)
		}
	}

	nabu.FromMessage("Node shutdown completed: " + x.Host.Name).Log()
	return nil
}
