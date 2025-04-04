package query

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEvaluateOperation_CoversAllOperationsRegistryTypes(t *testing.T) {
	// Dummy inputs by type
	dummies := map[string][2]any{
		"string":    {"a", "a"},
		"bool":      {true, true},
		"int":       {1, 1},
		"int8":      {int8(1), int8(1)},
		"int16":     {int16(1), int16(1)},
		"int32":     {int32(1), int32(1)},
		"int64":     {int64(1), int64(1)},
		"uint":      {uint(1), uint(1)},
		"uint8":     {uint8(1), uint8(1)},
		"uint16":    {uint16(1), uint16(1)},
		"uint32":    {uint32(1), uint32(1)},
		"uint64":    {uint64(1), uint64(1)},
		"float32":   {float32(1.0), float32(1.0)},
		"float64":   {float64(1.0), float64(1.0)},
		"uuid.UUID": {uuid.New(), uuid.New()},
		"time.Time": {time.Now(), time.Now()},
	}

	for fieldType := range OperatorsRegistry {
		d, ok := dummies[fieldType]
		if !ok {
			t.Errorf("missing test dummy values for type: %s", fieldType)
			continue
		}

		_, err := EvaluateOperation(OperatorTypeEqual, fieldType, d[0], d[1])
		if err != nil {
			t.Errorf("EvaluateOperation missing case for type %q or failed: %v", fieldType, err)
		}
	}
}

func TestStringOperations_All(t *testing.T) {
	tests := []struct {
		op       OperatorType
		a, b     string
		expected bool
	}{
		{OperatorTypeEqual, "hello", "hello", true},
		{OperatorTypeEqual, "hello", "world", false},
		{OperatorTypeNotEqual, "hello", "world", true},
		{OperatorTypeNotEqual, "hello", "hello", false},
		{OperatorTypeContains, "hello world", "world", true},
		{OperatorTypeContains, "hello world", "mars", false},
		{OperatorTypeNotContains, "hello world", "mars", true},
		{OperatorTypeNotContains, "hello world", "world", false},
		{OperatorTypeStartsWith, "golang", "go", true},
		{OperatorTypeStartsWith, "golang", "lang", false},
		{OperatorTypeEndsWith, "golang", "lang", true},
		{OperatorTypeEndsWith, "golang", "go", false},
	}

	for _, tt := range tests {
		if got := StringOperations[tt.op](tt.a, tt.b); got != tt.expected {
			t.Errorf("StringOperations[%v](%q, %q) = %v; want %v", tt.op, tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestBoolOperations(t *testing.T) {
	if !BoolOperations[OperatorTypeEqual](true, true) {
		t.Error("true == true should be true")
	}
	if BoolOperations[OperatorTypeEqual](true, false) {
		t.Error("true == false should be false")
	}
	if BoolOperations[OperatorTypeNotEqual](true, false) == false {
		t.Error("true != false should be true")
	}
	if BoolOperations[OperatorTypeNotEqual](true, true) {
		t.Error("true != true should be false")
	}
}

func TestIntOperations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int) bool
		a, b     int
		expected bool
	}{
		{"== true", IntOperations[OperatorTypeEqual], 5, 5, true},
		{"!= true", IntOperations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", IntOperations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", IntOperations[OperatorTypeLessThan], 4, 5, true},
		{">= true", IntOperations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", IntOperations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt8Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int8) bool
		a, b     int8
		expected bool
	}{
		{"== true", Int8Operations[OperatorTypeEqual], 1, 1, true},
		{"!= false", Int8Operations[OperatorTypeNotEqual], 1, 1, false},
		{"> true", Int8Operations[OperatorTypeGreaterThan], 2, 1, true},
		{"< true", Int8Operations[OperatorTypeLessThan], 1, 2, true},
		{">= true", Int8Operations[OperatorTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int8Operations[OperatorTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int8 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt16Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int16) bool
		a, b     int16
		expected bool
	}{
		{"== true", Int16Operations[OperatorTypeEqual], 1, 1, true},
		{"!= false", Int16Operations[OperatorTypeNotEqual], 1, 1, false},
		{"> true", Int16Operations[OperatorTypeGreaterThan], 2, 1, true},
		{"< true", Int16Operations[OperatorTypeLessThan], 1, 2, true},
		{">= true", Int16Operations[OperatorTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int16Operations[OperatorTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int16 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt32Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int32) bool
		a, b     int32
		expected bool
	}{
		{"== true", Int32Operations[OperatorTypeEqual], 1, 1, true},
		{"!= false", Int32Operations[OperatorTypeNotEqual], 1, 1, false},
		{"> true", Int32Operations[OperatorTypeGreaterThan], 2, 1, true},
		{"< true", Int32Operations[OperatorTypeLessThan], 1, 2, true},
		{">= true", Int32Operations[OperatorTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int32Operations[OperatorTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt64Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int64) bool
		a, b     int64
		expected bool
	}{
		{"== true", Int64Operations[OperatorTypeEqual], 1, 1, true},
		{"!= false", Int64Operations[OperatorTypeNotEqual], 1, 1, false},
		{"> true", Int64Operations[OperatorTypeGreaterThan], 2, 1, true},
		{"< true", Int64Operations[OperatorTypeLessThan], 1, 2, true},
		{">= true", Int64Operations[OperatorTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int64Operations[OperatorTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUintOperations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint) bool
		a, b     uint
		expected bool
	}{
		{"== true", UintOperations[OperatorTypeEqual], 5, 5, true},
		{"!= true", UintOperations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", UintOperations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", UintOperations[OperatorTypeLessThan], 4, 5, true},
		{">= true", UintOperations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", UintOperations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint8Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint8) bool
		a, b     uint8
		expected bool
	}{
		{"== true", Uint8Operations[OperatorTypeEqual], 5, 5, true},
		{"!= true", Uint8Operations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", Uint8Operations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", Uint8Operations[OperatorTypeLessThan], 4, 5, true},
		{">= true", Uint8Operations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint8Operations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint8 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint16Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint16) bool
		a, b     uint16
		expected bool
	}{
		{"== true", Uint16Operations[OperatorTypeEqual], 5, 5, true},
		{"!= true", Uint16Operations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", Uint16Operations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", Uint16Operations[OperatorTypeLessThan], 4, 5, true},
		{">= true", Uint16Operations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint16Operations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint16 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint32Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint32) bool
		a, b     uint32
		expected bool
	}{
		{"== true", Uint32Operations[OperatorTypeEqual], 5, 5, true},
		{"!= true", Uint32Operations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", Uint32Operations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", Uint32Operations[OperatorTypeLessThan], 4, 5, true},
		{">= true", Uint32Operations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint32Operations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint64Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint64) bool
		a, b     uint64
		expected bool
	}{
		{"== true", Uint64Operations[OperatorTypeEqual], 5, 5, true},
		{"!= true", Uint64Operations[OperatorTypeNotEqual], 5, 6, true},
		{"> true", Uint64Operations[OperatorTypeGreaterThan], 7, 6, true},
		{"< true", Uint64Operations[OperatorTypeLessThan], 4, 5, true},
		{">= true", Uint64Operations[OperatorTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint64Operations[OperatorTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestFloat32Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b float32) bool
		a, b     float32
		expected bool
	}{
		{"== true", Float32Operations[OperatorTypeEqual], 1.1, 1.1, true},
		{"!= true", Float32Operations[OperatorTypeNotEqual], 1.1, 2.2, true},
		{"> true", Float32Operations[OperatorTypeGreaterThan], 3.5, 2.5, true},
		{"< true", Float32Operations[OperatorTypeLessThan], 1.0, 2.0, true},
		{">= true", Float32Operations[OperatorTypeGreaterThanEqual], 4.0, 4.0, true},
		{"<= true", Float32Operations[OperatorTypeLessThanEqual], 2.2, 2.2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Float32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestFloat64Operations(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b float64) bool
		a, b     float64
		expected bool
	}{
		{"== true", Float64Operations[OperatorTypeEqual], 1.1, 1.1, true},
		{"!= true", Float64Operations[OperatorTypeNotEqual], 1.1, 2.2, true},
		{"> true", Float64Operations[OperatorTypeGreaterThan], 3.5, 2.5, true},
		{"< true", Float64Operations[OperatorTypeLessThan], 1.0, 2.0, true},
		{">= true", Float64Operations[OperatorTypeGreaterThanEqual], 4.0, 4.0, true},
		{"<= true", Float64Operations[OperatorTypeLessThanEqual], 2.2, 2.2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Float64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUuidOperations(t *testing.T) {
	a := uuid.New()
	b := uuid.New()
	cases := []struct {
		name     string
		op       func(a, b uuid.UUID) bool
		a, b     uuid.UUID
		expected bool
	}{
		{"== true", UuidOperations[OperatorTypeEqual], a, a, true},
		{"!= true", UuidOperations[OperatorTypeNotEqual], a, b, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("UUID %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestTimeOperations(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Minute)
	cases := []struct {
		name     string
		op       func(a, b time.Time) bool
		a, b     time.Time
		expected bool
	}{
		{"== true", TimeOperations[OperatorTypeEqual], now, now, true},
		{"!= true", TimeOperations[OperatorTypeNotEqual], now, later, true},
		{"> true", TimeOperations[OperatorTypeGreaterThan], later, now, true},
		{"< true", TimeOperations[OperatorTypeLessThan], now, later, true},
		{">= true", TimeOperations[OperatorTypeGreaterThanEqual], now, now, true},
		{"<= true", TimeOperations[OperatorTypeLessThanEqual], now, now, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Time %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}
