package main

import (
	"errors"
)

var (
	ErrPathConfigNotSpecified = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound = errors.New("pathConfig: specified path cannot be found")
)
