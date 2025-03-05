package gen

import (
	"errors"
)

var (
	ErrGeneratorStructNotFound         = errors.New("generator: struct not found")
	ErrGeneratorStructMoreThanOneFound = errors.New("generator: more than one struct found, there should be a single one")
)
