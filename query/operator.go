package query

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type OpType int

const (
	OpTypeUndefined OpType = iota

	OpTypeEqual
	OpTypeNotEqual
	OpTypeGreaterThan
	OpTypeGreaterThanEqual
	OpTypeLessThan
	OpTypeLessThanEqual
	OpTypeContains
	OpTypeNotContain
	OpTypeStartsWith
	OpTypeEndsWith
)

var UuidOps = map[OpType]func(a, b uuid.UUID) bool{
	OpTypeEqual:    func(a, b uuid.UUID) bool { return a == b },
	OpTypeNotEqual: func(a, b uuid.UUID) bool { return a != b },
}

var StringOps = map[OpType]func(a, b string) bool{
	OpTypeEqual:      func(a, b string) bool { return a == b },
	OpTypeNotEqual:   func(a, b string) bool { return a != b },
	OpTypeContains:   func(a, b string) bool { return strings.Contains(a, b) },
	OpTypeNotContain: func(a, b string) bool { return !strings.Contains(a, b) },
	OpTypeStartsWith: func(a, b string) bool { return strings.HasPrefix(a, b) },
	OpTypeEndsWith:   func(a, b string) bool { return strings.HasSuffix(a, b) },
}

var BoolOps = map[OpType]func(a, b bool) bool{
	OpTypeEqual:    func(a, b bool) bool { return a == b },
	OpTypeNotEqual: func(a, b bool) bool { return a != b },
}

var TimeOps = map[OpType]func(a, b time.Time) bool{
	OpTypeEqual:            func(a, b time.Time) bool { return a.Equal(b) },
	OpTypeNotEqual:         func(a, b time.Time) bool { return !a.Equal(b) },
	OpTypeGreaterThan:      func(a, b time.Time) bool { return a.After(b) },
	OpTypeLessThan:         func(a, b time.Time) bool { return a.Before(b) },
	OpTypeGreaterThanEqual: func(a, b time.Time) bool { return a.After(b) || a.Equal(b) },
	OpTypeLessThanEqual:    func(a, b time.Time) bool { return a.Before(b) || a.Equal(b) },
}

var IntOps = map[OpType]func(a, b int) bool{
	OpTypeEqual:            func(a, b int) bool { return a == b },
	OpTypeNotEqual:         func(a, b int) bool { return a != b },
	OpTypeGreaterThan:      func(a, b int) bool { return a > b },
	OpTypeLessThan:         func(a, b int) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b int) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b int) bool { return a <= b },
}
var Int8Ops = map[OpType]func(a, b int8) bool{
	OpTypeEqual:            func(a, b int8) bool { return a == b },
	OpTypeNotEqual:         func(a, b int8) bool { return a != b },
	OpTypeGreaterThan:      func(a, b int8) bool { return a > b },
	OpTypeLessThan:         func(a, b int8) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b int8) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b int8) bool { return a <= b },
}
var Int16Ops = map[OpType]func(a, b int16) bool{
	OpTypeEqual:            func(a, b int16) bool { return a == b },
	OpTypeNotEqual:         func(a, b int16) bool { return a != b },
	OpTypeGreaterThan:      func(a, b int16) bool { return a > b },
	OpTypeLessThan:         func(a, b int16) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b int16) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b int16) bool { return a <= b },
}
var Int32Ops = map[OpType]func(a, b int32) bool{
	OpTypeEqual:            func(a, b int32) bool { return a == b },
	OpTypeNotEqual:         func(a, b int32) bool { return a != b },
	OpTypeGreaterThan:      func(a, b int32) bool { return a > b },
	OpTypeLessThan:         func(a, b int32) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b int32) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b int32) bool { return a <= b },
}
var Int64Ops = map[OpType]func(a, b int64) bool{
	OpTypeEqual:            func(a, b int64) bool { return a == b },
	OpTypeNotEqual:         func(a, b int64) bool { return a != b },
	OpTypeGreaterThan:      func(a, b int64) bool { return a > b },
	OpTypeLessThan:         func(a, b int64) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b int64) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b int64) bool { return a <= b },
}

var UintOps = map[OpType]func(a, b uint) bool{
	OpTypeEqual:            func(a, b uint) bool { return a == b },
	OpTypeNotEqual:         func(a, b uint) bool { return a != b },
	OpTypeGreaterThan:      func(a, b uint) bool { return a > b },
	OpTypeLessThan:         func(a, b uint) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b uint) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b uint) bool { return a <= b },
}
var Uint8Ops = map[OpType]func(a, b uint8) bool{
	OpTypeEqual:            func(a, b uint8) bool { return a == b },
	OpTypeNotEqual:         func(a, b uint8) bool { return a != b },
	OpTypeGreaterThan:      func(a, b uint8) bool { return a > b },
	OpTypeLessThan:         func(a, b uint8) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b uint8) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b uint8) bool { return a <= b },
}
var Uint16Ops = map[OpType]func(a, b uint16) bool{
	OpTypeEqual:            func(a, b uint16) bool { return a == b },
	OpTypeNotEqual:         func(a, b uint16) bool { return a != b },
	OpTypeGreaterThan:      func(a, b uint16) bool { return a > b },
	OpTypeLessThan:         func(a, b uint16) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b uint16) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b uint16) bool { return a <= b },
}
var Uint32Ops = map[OpType]func(a, b uint32) bool{
	OpTypeEqual:            func(a, b uint32) bool { return a == b },
	OpTypeNotEqual:         func(a, b uint32) bool { return a != b },
	OpTypeGreaterThan:      func(a, b uint32) bool { return a > b },
	OpTypeLessThan:         func(a, b uint32) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b uint32) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b uint32) bool { return a <= b },
}
var Uint64Ops = map[OpType]func(a, b uint64) bool{
	OpTypeEqual:            func(a, b uint64) bool { return a == b },
	OpTypeNotEqual:         func(a, b uint64) bool { return a != b },
	OpTypeGreaterThan:      func(a, b uint64) bool { return a > b },
	OpTypeLessThan:         func(a, b uint64) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b uint64) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b uint64) bool { return a <= b },
}

var Float32Ops = map[OpType]func(a, b float32) bool{
	OpTypeEqual:            func(a, b float32) bool { return a == b },
	OpTypeNotEqual:         func(a, b float32) bool { return a != b },
	OpTypeGreaterThan:      func(a, b float32) bool { return a > b },
	OpTypeLessThan:         func(a, b float32) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b float32) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b float32) bool { return a <= b },
}
var Float64Ops = map[OpType]func(a, b float64) bool{
	OpTypeEqual:            func(a, b float64) bool { return a == b },
	OpTypeNotEqual:         func(a, b float64) bool { return a != b },
	OpTypeGreaterThan:      func(a, b float64) bool { return a > b },
	OpTypeLessThan:         func(a, b float64) bool { return a < b },
	OpTypeGreaterThanEqual: func(a, b float64) bool { return a >= b },
	OpTypeLessThanEqual:    func(a, b float64) bool { return a <= b },
}
