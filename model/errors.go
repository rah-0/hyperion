package model

import (
	"errors"
)

var (
	ErrGeneratorStructNotFound = errors.New("generator: struct not found")

	ErrMessageEmpty = errors.New("message: is empty")

	ErrConfigNodesNotFound       = errors.New("GlobalConfig: node list is empty")
	ErrConfigNodeNotFoundForHost = errors.New("GlobalConfig: node not found for current hostname")

	ErrPathConfigNoContent    = errors.New("pathConfig: GlobalConfig file is empty")
	ErrPathConfigNotSpecified = errors.New("pathConfig: not specified in either command line argument or environment variable")
	ErrPathConfigNotFound     = errors.New("pathConfig: specified path cannot be found")
	ErrPathConfigNotEditable  = errors.New("pathConfig: specified path is not editable")

	ErrQueryNil                         = errors.New("query: cannot be nil")
	ErrQueryEntityNoUuid                = errors.New("query: entity has no uuid")
	ErrQueryEntityFieldNotFound         = errors.New("query: entity field not found")
	ErrQueryEntityFieldOperatorNotFound = errors.New("query: operator not found for given field")
)
