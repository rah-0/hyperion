package query

type FilterType int

const (
	FilterTypeUndefined FilterType = iota
	FilterTypeOr
	FilterTypeAnd
)

type Filter struct {
	Field int
	Value any
	Op    OpType
}

type Filters struct {
	Type    FilterType
	Filters []Filter
}

type Query struct {
	Filters Filters
	Limit   int
}

func NewQuery() *Query {
	return &Query{}
}

func (x *Query) SetFilters(filters Filters) *Query {
	x.Filters = filters
	return x
}

func (x *Query) SetLimit(limit int) *Query {
	x.Limit = limit
	return x
}
