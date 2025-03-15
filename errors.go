package main

import (
	"errors"
)

var (
	ErrGeneratorStructNotFound = errors.New("generator: struct not found")

	ErrMessageEmpty = errors.New("message: is empty")

	ErrConfigNodesNotFound       = errors.New("GlobalConfig: node list is empty")
	ErrConfigNodeNotFoundForHost = errors.New("GlobalConfig: node not found for current hostname")

	ErrConnectionReachedRetryLimit = errors.New("Connection: retry limit reached")

	ErrPathConfigNoContent    = errors.New("pathConfig: GlobalConfig file is empty")
	ErrPathConfigNotSpecified = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound     = errors.New("pathConfig: specified path cannot be found")
	ErrPathConfigNotEditable  = errors.New("pathConfig: specified path is not editable")
)
