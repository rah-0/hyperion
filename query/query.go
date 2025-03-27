package query

import (
	"sort"
)

type Filter struct {
	Field int
	Value any
	Op    OpType
}

type Query struct {
	Filters []Filter
}

func NewQuery() *Query {
	return &Query{
		Filters: []Filter{},
	}
}

func (x *Query) AddFilter(field int, op OpType, value any) *Query {
	x.Filters = append(x.Filters, Filter{
		Field: field,
		Op:    op,
		Value: value,
	})
	return x
}

// SortFilters orders filters by ascending Field ID.
func (x *Query) SortFilters() {
	sort.SliceStable(x.Filters, func(i, j int) bool {
		return x.Filters[i].Field < x.Filters[j].Field
	})
}
