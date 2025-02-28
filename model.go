package main

import (
	"errors"
	"sync"
)

var (
	ErrConfigNodePortStatusNotAvailable = errors.New("node: port for status already in use")
	ErrConfigCurrentNodesNotFound       = errors.New("config: nodes not found for current hostname")
	ErrConfigNoNodes                    = errors.New("config: node list is empty")
	ErrPathConfigNoContent              = errors.New("pathConfig: config file is empty")
	ErrPathConfigNotSpecified           = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound               = errors.New("pathConfig: specified path cannot be found")
	ErrPathConfigNotEditable            = errors.New("pathConfig: specified path is not editable")
)

type Ports struct {
	Status int
	Data   int
}

type Path struct {
	Data string // Where data will be stored
}

type Host struct {
	Name  string
	Ports Ports
}

type Node struct {
	Host   Host
	Path   Path
	errCh  chan error
	Status NodeStatus

	mu sync.Mutex
}

type Config struct {
	ClusterName string
	Nodes       []Node
}
