package query

import (
	"fmt"
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
	OpTypeNotContains
	OpTypeStartsWith
	OpTypeEndsWith
)

var OpsRegistry = map[string]any{
	"string": StringOps,

	"bool": BoolOps,

	"int":   IntOps,
	"int8":  Int8Ops,
	"int16": Int16Ops,
	"int32": Int32Ops,
	"int64": Int64Ops,

	"uint":   UintOps,
	"uint8":  Uint8Ops,
	"uint16": Uint16Ops,
	"uint32": Uint32Ops,
	"uint64": Uint64Ops,

	"float32": Float32Ops,
	"float64": Float64Ops,

	"uuid.UUID": UuidOps,

	"time.Time": TimeOps,
}

func EvaluateOp(op OpType, fieldType string, a, b any) (bool, error) {
	switch fieldType {
	case "string":
		av, aok := a.(string)
		bv, bok := b.(string)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for string")
		}
		fn, ok := StringOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for string")
		}
		return fn(av, bv), nil

	case "bool":
		av, aok := a.(bool)
		bv, bok := b.(bool)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for bool")
		}
		fn, ok := BoolOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for bool")
		}
		return fn(av, bv), nil

	case "int":
		av, aok := a.(int)
		bv, bok := b.(int)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int")
		}
		fn, ok := IntOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for int")
		}
		return fn(av, bv), nil

	case "int8":
		av, aok := a.(int8)
		bv, bok := b.(int8)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int8")
		}
		fn, ok := Int8Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for int8")
		}
		return fn(av, bv), nil

	case "int16":
		av, aok := a.(int16)
		bv, bok := b.(int16)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int16")
		}
		fn, ok := Int16Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for int16")
		}
		return fn(av, bv), nil

	case "int32":
		av, aok := a.(int32)
		bv, bok := b.(int32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int32")
		}
		fn, ok := Int32Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for int32")
		}
		return fn(av, bv), nil

	case "int64":
		av, aok := a.(int64)
		bv, bok := b.(int64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int64")
		}
		fn, ok := Int64Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for int64")
		}
		return fn(av, bv), nil

	case "uint":
		av, aok := a.(uint)
		bv, bok := b.(uint)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint")
		}
		fn, ok := UintOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uint")
		}
		return fn(av, bv), nil

	case "uint8":
		av, aok := a.(uint8)
		bv, bok := b.(uint8)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint8")
		}
		fn, ok := Uint8Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uint8")
		}
		return fn(av, bv), nil

	case "uint16":
		av, aok := a.(uint16)
		bv, bok := b.(uint16)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint16")
		}
		fn, ok := Uint16Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uint16")
		}
		return fn(av, bv), nil

	case "uint32":
		av, aok := a.(uint32)
		bv, bok := b.(uint32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint32")
		}
		fn, ok := Uint32Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uint32")
		}
		return fn(av, bv), nil

	case "uint64":
		av, aok := a.(uint64)
		bv, bok := b.(uint64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint64")
		}
		fn, ok := Uint64Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uint64")
		}
		return fn(av, bv), nil

	case "float32":
		av, aok := a.(float32)
		bv, bok := b.(float32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for float32")
		}
		fn, ok := Float32Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for float32")
		}
		return fn(av, bv), nil

	case "float64":
		av, aok := a.(float64)
		bv, bok := b.(float64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for float64")
		}
		fn, ok := Float64Ops[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for float64")
		}
		return fn(av, bv), nil

	case "uuid.UUID":
		av, aok := a.(uuid.UUID)
		bv, bok := b.(uuid.UUID)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uuid.UUID")
		}
		fn, ok := UuidOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for uuid.UUID")
		}
		return fn(av, bv), nil

	case "time.Time":
		av, aok := a.(time.Time)
		bv, bok := b.(time.Time)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for time.Time")
		}
		fn, ok := TimeOps[op]
		if !ok {
			return false, fmt.Errorf("unsupported op for time.Time")
		}
		return fn(av, bv), nil

	default:
		return false, fmt.Errorf("unsupported field type: %s", fieldType)
	}
}

var StringOps = map[OpType]func(a, b string) bool{
	OpTypeEqual:       func(a, b string) bool { return a == b },
	OpTypeNotEqual:    func(a, b string) bool { return a != b },
	OpTypeContains:    func(a, b string) bool { return strings.Contains(a, b) },
	OpTypeNotContains: func(a, b string) bool { return !strings.Contains(a, b) },
	OpTypeStartsWith:  func(a, b string) bool { return strings.HasPrefix(a, b) },
	OpTypeEndsWith:    func(a, b string) bool { return strings.HasSuffix(a, b) },
}

var BoolOps = map[OpType]func(a, b bool) bool{
	OpTypeEqual:    func(a, b bool) bool { return a == b },
	OpTypeNotEqual: func(a, b bool) bool { return a != b },
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

var UuidOps = map[OpType]func(a, b uuid.UUID) bool{
	OpTypeEqual:    func(a, b uuid.UUID) bool { return a == b },
	OpTypeNotEqual: func(a, b uuid.UUID) bool { return a != b },
}

var TimeOps = map[OpType]func(a, b time.Time) bool{
	OpTypeEqual:            func(a, b time.Time) bool { return a.Equal(b) },
	OpTypeNotEqual:         func(a, b time.Time) bool { return !a.Equal(b) },
	OpTypeGreaterThan:      func(a, b time.Time) bool { return a.After(b) },
	OpTypeLessThan:         func(a, b time.Time) bool { return a.Before(b) },
	OpTypeGreaterThanEqual: func(a, b time.Time) bool { return a.After(b) || a.Equal(b) },
	OpTypeLessThanEqual:    func(a, b time.Time) bool { return a.Before(b) || a.Equal(b) },
}
