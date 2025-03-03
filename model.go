package main

import (
	"errors"
)

var (
	ErrMessageEmpty = errors.New("message: is empty")

	ErrNodePortNotAvailable = errors.New("node: port not available")

	ErrConfigNodesNotFound        = errors.New("config: node list is empty")
	ErrConfigNodesNotFoundForHost = errors.New("config: nodes not found for current hostname")

	ErrPathConfigNoContent    = errors.New("pathConfig: config file is empty")
	ErrPathConfigNotSpecified = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound     = errors.New("pathConfig: specified path cannot be found")
	ErrPathConfigNotEditable  = errors.New("pathConfig: specified path is not editable")
)

type Path struct {
	Data string // Where data will be stored
}

type Host struct {
	Name string
	Port int
}

type Config struct {
	ClusterName string
	Nodes       []*Node
}
