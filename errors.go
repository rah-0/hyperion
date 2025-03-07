package main

import (
	"errors"
)

var (
	ErrGeneratorStructNotFound = errors.New("generator: struct not found")

	ErrMessageEmpty = errors.New("message: is empty")

	ErrConfigNodesNotFound        = errors.New("config: node list is empty")
	ErrConfigNodesNotFoundForHost = errors.New("config: nodes not found for current hostname")

	ErrPathConfigNoContent    = errors.New("pathConfig: config file is empty")
	ErrPathConfigNotSpecified = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound     = errors.New("pathConfig: specified path cannot be found")
	ErrPathConfigNotEditable  = errors.New("pathConfig: specified path is not editable")
)
