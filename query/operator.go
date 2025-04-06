package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OperatorType int

const (
	OperatorTypeUndefined OperatorType = iota

	OperatorTypeEqual
	OperatorTypeNotEqual
	OperatorTypeGreater
	OperatorTypeGreaterOrEqual
	OperatorTypeLesser
	OperatorTypeLesserOrEqual
	OperatorTypeContains
	OperatorTypeNotContains
	OperatorTypeStartsWith
	OperatorTypeEndsWith
)

var OperatorsRegistry = map[string]any{
	"string": StringOperations,

	"bool": BoolOperations,

	"int":   IntOperations,
	"int8":  Int8Operations,
	"int16": Int16Operations,
	"int32": Int32Operations,
	"int64": Int64Operations,

	"uint":   UintOperations,
	"uint8":  Uint8Operations,
	"uint16": Uint16Operations,
	"uint32": Uint32Operations,
	"uint64": Uint64Operations,

	"float32": Float32Operations,
	"float64": Float64Operations,

	"uuid.UUID": UuidOperations,

	"time.Time": TimeOperations,
}

func EvaluateOperation(operator OperatorType, fieldType string, a, b any) (bool, error) {
	switch fieldType {
	case "string":
		av, aok := a.(string)
		bv, bok := b.(string)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for string")
		}
		fn, ok := StringOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for string")
		}
		return fn(av, bv), nil

	case "bool":
		av, aok := a.(bool)
		bv, bok := b.(bool)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for bool")
		}
		fn, ok := BoolOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for bool")
		}
		return fn(av, bv), nil

	case "int":
		av, aok := a.(int)
		bv, bok := b.(int)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int")
		}
		fn, ok := IntOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for int")
		}
		return fn(av, bv), nil

	case "int8":
		av, aok := a.(int8)
		bv, bok := b.(int8)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int8")
		}
		fn, ok := Int8Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for int8")
		}
		return fn(av, bv), nil

	case "int16":
		av, aok := a.(int16)
		bv, bok := b.(int16)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int16")
		}
		fn, ok := Int16Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for int16")
		}
		return fn(av, bv), nil

	case "int32":
		av, aok := a.(int32)
		bv, bok := b.(int32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int32")
		}
		fn, ok := Int32Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for int32")
		}
		return fn(av, bv), nil

	case "int64":
		av, aok := a.(int64)
		bv, bok := b.(int64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for int64")
		}
		fn, ok := Int64Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for int64")
		}
		return fn(av, bv), nil

	case "uint":
		av, aok := a.(uint)
		bv, bok := b.(uint)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint")
		}
		fn, ok := UintOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uint")
		}
		return fn(av, bv), nil

	case "uint8":
		av, aok := a.(uint8)
		bv, bok := b.(uint8)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint8")
		}
		fn, ok := Uint8Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uint8")
		}
		return fn(av, bv), nil

	case "uint16":
		av, aok := a.(uint16)
		bv, bok := b.(uint16)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint16")
		}
		fn, ok := Uint16Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uint16")
		}
		return fn(av, bv), nil

	case "uint32":
		av, aok := a.(uint32)
		bv, bok := b.(uint32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint32")
		}
		fn, ok := Uint32Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uint32")
		}
		return fn(av, bv), nil

	case "uint64":
		av, aok := a.(uint64)
		bv, bok := b.(uint64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uint64")
		}
		fn, ok := Uint64Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uint64")
		}
		return fn(av, bv), nil

	case "float32":
		av, aok := a.(float32)
		bv, bok := b.(float32)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for float32")
		}
		fn, ok := Float32Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for float32")
		}
		return fn(av, bv), nil

	case "float64":
		av, aok := a.(float64)
		bv, bok := b.(float64)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for float64")
		}
		fn, ok := Float64Operations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for float64")
		}
		return fn(av, bv), nil

	case "uuid.UUID":
		av, aok := a.(uuid.UUID)
		bv, bok := b.(uuid.UUID)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for uuid.UUID")
		}
		fn, ok := UuidOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for uuid.UUID")
		}
		return fn(av, bv), nil

	case "time.Time":
		av, aok := a.(time.Time)
		bv, bok := b.(time.Time)
		if !aok || !bok {
			return false, fmt.Errorf("type mismatch for time.Time")
		}
		fn, ok := TimeOperations[operator]
		if !ok {
			return false, fmt.Errorf("unsupported operator for time.Time")
		}
		return fn(av, bv), nil

	default:
		return false, fmt.Errorf("unsupported field type: %s", fieldType)
	}
}

var StringOperations = map[OperatorType]func(a, b string) bool{
	OperatorTypeEqual:       func(a, b string) bool { return a == b },
	OperatorTypeNotEqual:    func(a, b string) bool { return a != b },
	OperatorTypeContains:    func(a, b string) bool { return strings.Contains(a, b) },
	OperatorTypeNotContains: func(a, b string) bool { return !strings.Contains(a, b) },
	OperatorTypeStartsWith:  func(a, b string) bool { return strings.HasPrefix(a, b) },
	OperatorTypeEndsWith:    func(a, b string) bool { return strings.HasSuffix(a, b) },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b string) bool { return a > b },
	OperatorTypeLesser:         func(a, b string) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b string) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b string) bool { return a <= b },
}

var BoolOperations = map[OperatorType]func(a, b bool) bool{
	OperatorTypeEqual:    func(a, b bool) bool { return a == b },
	OperatorTypeNotEqual: func(a, b bool) bool { return a != b },
}

var IntOperations = map[OperatorType]func(a, b int) bool{
	OperatorTypeEqual:    func(a, b int) bool { return a == b },
	OperatorTypeNotEqual: func(a, b int) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b int) bool { return a > b },
	OperatorTypeLesser:         func(a, b int) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b int) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b int) bool { return a <= b },
}
var Int8Operations = map[OperatorType]func(a, b int8) bool{
	OperatorTypeEqual:    func(a, b int8) bool { return a == b },
	OperatorTypeNotEqual: func(a, b int8) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b int8) bool { return a > b },
	OperatorTypeLesser:         func(a, b int8) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b int8) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b int8) bool { return a <= b },
}
var Int16Operations = map[OperatorType]func(a, b int16) bool{
	OperatorTypeEqual:    func(a, b int16) bool { return a == b },
	OperatorTypeNotEqual: func(a, b int16) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b int16) bool { return a > b },
	OperatorTypeLesser:         func(a, b int16) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b int16) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b int16) bool { return a <= b },
}
var Int32Operations = map[OperatorType]func(a, b int32) bool{
	OperatorTypeEqual:    func(a, b int32) bool { return a == b },
	OperatorTypeNotEqual: func(a, b int32) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b int32) bool { return a > b },
	OperatorTypeLesser:         func(a, b int32) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b int32) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b int32) bool { return a <= b },
}
var Int64Operations = map[OperatorType]func(a, b int64) bool{
	OperatorTypeEqual:    func(a, b int64) bool { return a == b },
	OperatorTypeNotEqual: func(a, b int64) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b int64) bool { return a > b },
	OperatorTypeLesser:         func(a, b int64) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b int64) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b int64) bool { return a <= b },
}

var UintOperations = map[OperatorType]func(a, b uint) bool{
	OperatorTypeEqual:    func(a, b uint) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uint) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b uint) bool { return a > b },
	OperatorTypeLesser:         func(a, b uint) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b uint) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b uint) bool { return a <= b },
}
var Uint8Operations = map[OperatorType]func(a, b uint8) bool{
	OperatorTypeEqual:    func(a, b uint8) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uint8) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b uint8) bool { return a > b },
	OperatorTypeLesser:         func(a, b uint8) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b uint8) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b uint8) bool { return a <= b },
}
var Uint16Operations = map[OperatorType]func(a, b uint16) bool{
	OperatorTypeEqual:    func(a, b uint16) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uint16) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b uint16) bool { return a > b },
	OperatorTypeLesser:         func(a, b uint16) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b uint16) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b uint16) bool { return a <= b },
}
var Uint32Operations = map[OperatorType]func(a, b uint32) bool{
	OperatorTypeEqual:    func(a, b uint32) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uint32) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b uint32) bool { return a > b },
	OperatorTypeLesser:         func(a, b uint32) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b uint32) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b uint32) bool { return a <= b },
}
var Uint64Operations = map[OperatorType]func(a, b uint64) bool{
	OperatorTypeEqual:    func(a, b uint64) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uint64) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b uint64) bool { return a > b },
	OperatorTypeLesser:         func(a, b uint64) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b uint64) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b uint64) bool { return a <= b },
}

var Float32Operations = map[OperatorType]func(a, b float32) bool{
	OperatorTypeEqual:    func(a, b float32) bool { return a == b },
	OperatorTypeNotEqual: func(a, b float32) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b float32) bool { return a > b },
	OperatorTypeLesser:         func(a, b float32) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b float32) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b float32) bool { return a <= b },
}
var Float64Operations = map[OperatorType]func(a, b float64) bool{
	OperatorTypeEqual:    func(a, b float64) bool { return a == b },
	OperatorTypeNotEqual: func(a, b float64) bool { return a != b },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b float64) bool { return a > b },
	OperatorTypeLesser:         func(a, b float64) bool { return a < b },
	OperatorTypeGreaterOrEqual: func(a, b float64) bool { return a >= b },
	OperatorTypeLesserOrEqual:  func(a, b float64) bool { return a <= b },
}

var UuidOperations = map[OperatorType]func(a, b uuid.UUID) bool{
	OperatorTypeEqual:    func(a, b uuid.UUID) bool { return a == b },
	OperatorTypeNotEqual: func(a, b uuid.UUID) bool { return a != b },
}

var TimeOperations = map[OperatorType]func(a, b time.Time) bool{
	OperatorTypeEqual:    func(a, b time.Time) bool { return a.Equal(b) },
	OperatorTypeNotEqual: func(a, b time.Time) bool { return !a.Equal(b) },
	// Operators below also used for sorting
	OperatorTypeGreater:        func(a, b time.Time) bool { return a.After(b) },
	OperatorTypeLesser:         func(a, b time.Time) bool { return a.Before(b) },
	OperatorTypeGreaterOrEqual: func(a, b time.Time) bool { return a.After(b) || a.Equal(b) },
	OperatorTypeLesserOrEqual:  func(a, b time.Time) bool { return a.Before(b) || a.Equal(b) },
}
