package query

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEvaluateOp_CoversAllOpsRegistryTypes(t *testing.T) {
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

	for fieldType := range OpsRegistry {
		d, ok := dummies[fieldType]
		if !ok {
			t.Errorf("missing test dummy values for type: %s", fieldType)
			continue
		}

		_, err := EvaluateOp(OpTypeEqual, fieldType, d[0], d[1])
		if err != nil {
			t.Errorf("EvaluateOp missing case for type %q or failed: %v", fieldType, err)
		}
	}
}

func TestStringOps_All(t *testing.T) {
	tests := []struct {
		op       OpType
		a, b     string
		expected bool
	}{
		{OpTypeEqual, "hello", "hello", true},
		{OpTypeEqual, "hello", "world", false},
		{OpTypeNotEqual, "hello", "world", true},
		{OpTypeNotEqual, "hello", "hello", false},
		{OpTypeContains, "hello world", "world", true},
		{OpTypeContains, "hello world", "mars", false},
		{OpTypeNotContains, "hello world", "mars", true},
		{OpTypeNotContains, "hello world", "world", false},
		{OpTypeStartsWith, "golang", "go", true},
		{OpTypeStartsWith, "golang", "lang", false},
		{OpTypeEndsWith, "golang", "lang", true},
		{OpTypeEndsWith, "golang", "go", false},
	}

	for _, tt := range tests {
		if got := StringOps[tt.op](tt.a, tt.b); got != tt.expected {
			t.Errorf("StringOps[%v](%q, %q) = %v; want %v", tt.op, tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestBoolOps(t *testing.T) {
	if !BoolOps[OpTypeEqual](true, true) {
		t.Error("true == true should be true")
	}
	if BoolOps[OpTypeEqual](true, false) {
		t.Error("true == false should be false")
	}
	if BoolOps[OpTypeNotEqual](true, false) == false {
		t.Error("true != false should be true")
	}
	if BoolOps[OpTypeNotEqual](true, true) {
		t.Error("true != true should be false")
	}
}

func TestIntOps(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int) bool
		a, b     int
		expected bool
	}{
		{"== true", IntOps[OpTypeEqual], 5, 5, true},
		{"!= true", IntOps[OpTypeNotEqual], 5, 6, true},
		{"> true", IntOps[OpTypeGreaterThan], 7, 6, true},
		{"< true", IntOps[OpTypeLessThan], 4, 5, true},
		{">= true", IntOps[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", IntOps[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt8Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int8) bool
		a, b     int8
		expected bool
	}{
		{"== true", Int8Ops[OpTypeEqual], 1, 1, true},
		{"!= false", Int8Ops[OpTypeNotEqual], 1, 1, false},
		{"> true", Int8Ops[OpTypeGreaterThan], 2, 1, true},
		{"< true", Int8Ops[OpTypeLessThan], 1, 2, true},
		{">= true", Int8Ops[OpTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int8Ops[OpTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int8 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt16Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int16) bool
		a, b     int16
		expected bool
	}{
		{"== true", Int16Ops[OpTypeEqual], 1, 1, true},
		{"!= false", Int16Ops[OpTypeNotEqual], 1, 1, false},
		{"> true", Int16Ops[OpTypeGreaterThan], 2, 1, true},
		{"< true", Int16Ops[OpTypeLessThan], 1, 2, true},
		{">= true", Int16Ops[OpTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int16Ops[OpTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int16 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt32Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int32) bool
		a, b     int32
		expected bool
	}{
		{"== true", Int32Ops[OpTypeEqual], 1, 1, true},
		{"!= false", Int32Ops[OpTypeNotEqual], 1, 1, false},
		{"> true", Int32Ops[OpTypeGreaterThan], 2, 1, true},
		{"< true", Int32Ops[OpTypeLessThan], 1, 2, true},
		{">= true", Int32Ops[OpTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int32Ops[OpTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestInt64Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b int64) bool
		a, b     int64
		expected bool
	}{
		{"== true", Int64Ops[OpTypeEqual], 1, 1, true},
		{"!= false", Int64Ops[OpTypeNotEqual], 1, 1, false},
		{"> true", Int64Ops[OpTypeGreaterThan], 2, 1, true},
		{"< true", Int64Ops[OpTypeLessThan], 1, 2, true},
		{">= true", Int64Ops[OpTypeGreaterThanEqual], 2, 2, true},
		{"<= true", Int64Ops[OpTypeLessThanEqual], 2, 2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Int64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUintOps(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint) bool
		a, b     uint
		expected bool
	}{
		{"== true", UintOps[OpTypeEqual], 5, 5, true},
		{"!= true", UintOps[OpTypeNotEqual], 5, 6, true},
		{"> true", UintOps[OpTypeGreaterThan], 7, 6, true},
		{"< true", UintOps[OpTypeLessThan], 4, 5, true},
		{">= true", UintOps[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", UintOps[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint8Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint8) bool
		a, b     uint8
		expected bool
	}{
		{"== true", Uint8Ops[OpTypeEqual], 5, 5, true},
		{"!= true", Uint8Ops[OpTypeNotEqual], 5, 6, true},
		{"> true", Uint8Ops[OpTypeGreaterThan], 7, 6, true},
		{"< true", Uint8Ops[OpTypeLessThan], 4, 5, true},
		{">= true", Uint8Ops[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint8Ops[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint8 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint16Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint16) bool
		a, b     uint16
		expected bool
	}{
		{"== true", Uint16Ops[OpTypeEqual], 5, 5, true},
		{"!= true", Uint16Ops[OpTypeNotEqual], 5, 6, true},
		{"> true", Uint16Ops[OpTypeGreaterThan], 7, 6, true},
		{"< true", Uint16Ops[OpTypeLessThan], 4, 5, true},
		{">= true", Uint16Ops[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint16Ops[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint16 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint32Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint32) bool
		a, b     uint32
		expected bool
	}{
		{"== true", Uint32Ops[OpTypeEqual], 5, 5, true},
		{"!= true", Uint32Ops[OpTypeNotEqual], 5, 6, true},
		{"> true", Uint32Ops[OpTypeGreaterThan], 7, 6, true},
		{"< true", Uint32Ops[OpTypeLessThan], 4, 5, true},
		{">= true", Uint32Ops[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint32Ops[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUint64Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b uint64) bool
		a, b     uint64
		expected bool
	}{
		{"== true", Uint64Ops[OpTypeEqual], 5, 5, true},
		{"!= true", Uint64Ops[OpTypeNotEqual], 5, 6, true},
		{"> true", Uint64Ops[OpTypeGreaterThan], 7, 6, true},
		{"< true", Uint64Ops[OpTypeLessThan], 4, 5, true},
		{">= true", Uint64Ops[OpTypeGreaterThanEqual], 5, 5, true},
		{"<= true", Uint64Ops[OpTypeLessThanEqual], 5, 5, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Uint64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestFloat32Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b float32) bool
		a, b     float32
		expected bool
	}{
		{"== true", Float32Ops[OpTypeEqual], 1.1, 1.1, true},
		{"!= true", Float32Ops[OpTypeNotEqual], 1.1, 2.2, true},
		{"> true", Float32Ops[OpTypeGreaterThan], 3.5, 2.5, true},
		{"< true", Float32Ops[OpTypeLessThan], 1.0, 2.0, true},
		{">= true", Float32Ops[OpTypeGreaterThanEqual], 4.0, 4.0, true},
		{"<= true", Float32Ops[OpTypeLessThanEqual], 2.2, 2.2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Float32 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestFloat64Ops(t *testing.T) {
	cases := []struct {
		name     string
		op       func(a, b float64) bool
		a, b     float64
		expected bool
	}{
		{"== true", Float64Ops[OpTypeEqual], 1.1, 1.1, true},
		{"!= true", Float64Ops[OpTypeNotEqual], 1.1, 2.2, true},
		{"> true", Float64Ops[OpTypeGreaterThan], 3.5, 2.5, true},
		{"< true", Float64Ops[OpTypeLessThan], 1.0, 2.0, true},
		{">= true", Float64Ops[OpTypeGreaterThanEqual], 4.0, 4.0, true},
		{"<= true", Float64Ops[OpTypeLessThanEqual], 2.2, 2.2, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Float64 %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestUuidOps(t *testing.T) {
	a := uuid.New()
	b := uuid.New()
	cases := []struct {
		name     string
		op       func(a, b uuid.UUID) bool
		a, b     uuid.UUID
		expected bool
	}{
		{"== true", UuidOps[OpTypeEqual], a, a, true},
		{"!= true", UuidOps[OpTypeNotEqual], a, b, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("UUID %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}

func TestTimeOps(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Minute)
	cases := []struct {
		name     string
		op       func(a, b time.Time) bool
		a, b     time.Time
		expected bool
	}{
		{"== true", TimeOps[OpTypeEqual], now, now, true},
		{"!= true", TimeOps[OpTypeNotEqual], now, later, true},
		{"> true", TimeOps[OpTypeGreaterThan], later, now, true},
		{"< true", TimeOps[OpTypeLessThan], now, later, true},
		{">= true", TimeOps[OpTypeGreaterThanEqual], now, now, true},
		{"<= true", TimeOps[OpTypeLessThanEqual], now, now, true},
	}
	for _, c := range cases {
		if got := c.op(c.a, c.b); got != c.expected {
			t.Errorf("Time %s: got %v, want %v", c.name, got, c.expected)
		}
	}
}
