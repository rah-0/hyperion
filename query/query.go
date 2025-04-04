package query

type FilterType int
type OrderType int

const (
	FilterTypeUndefined FilterType = iota
	FilterTypeOr
	FilterTypeAnd
)
const (
	OrderTypeUndefined OrderType = iota
	OrderTypeAsc
	OrderTypeDesc
)

type Query struct {
	Filters Filters
	Orders  []Order
	Limit   int
}

type Filters struct {
	Type    FilterType
	Filters []Filter
}

type Filter struct {
	Field int
	Value any
	Op    OperatorType
}

type Order struct {
	Type  OrderType
	Field int
}

func NewQuery() *Query {
	return &Query{}
}

func (x *Query) SetFilters(opType FilterType, filters []Filter) *Query {
	x.Filters.Type = opType
	x.Filters.Filters = append(x.Filters.Filters, filters...)
	return x
}

func (x *Query) SetOrders(orders []Order) *Query {
	x.Orders = orders
	return x
}

func (x *Query) AddOrder(orderType OrderType, field int) *Query {
	x.Orders = append(x.Orders, Order{
		Type:  orderType,
		Field: field,
	})
	return x
}

func (x *Query) SetLimit(limit int) *Query {
	x.Limit = limit
	return x
}
