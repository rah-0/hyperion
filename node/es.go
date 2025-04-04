package node

import (
	"github.com/rah-0/parsort"

	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
)

func (x *EntityStorage) HandleQuery(q *query.Query) ([]register.Model, error) {
	mem := x.Memory.New().MemoryGetAll()
	fieldTypes := x.Memory.FieldTypes
	if q == nil {
		return nil, model.ErrQueryNil
	}

	hasOrders := len(q.Orders) > 0
	if hasOrders {
		if err := checkOrders(q, fieldTypes); err != nil {
			return nil, err
		}
		order(mem, q.Orders, fieldTypes)
	}

	hasFilters := q.Filters.Type != query.FilterTypeUndefined
	filters := q.Filters.Filters
	filterType := q.Filters.Type

	hasLimit := q.Limit > 0
	limit := q.Limit

	var results []register.Model
	for _, m := range mem {
		if !hasFilters || matchModel(m, filters, fieldTypes, filterType) {
			results = append(results, m)
			if hasLimit && len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

func matchModel(m register.Model, filters []query.Filter, fieldTypes map[int]string, qft query.FilterType) bool {
	switch qft {
	case query.FilterTypeOr:
		for _, f := range filters {
			a := m.GetFieldValue(f.Field)
			b := f.Value
			ok, _ := query.EvaluateOperation(f.Op, fieldTypes[f.Field], a, b)
			if ok {
				return true
			}
		}
		return false

	case query.FilterTypeAnd:
		for _, f := range filters {
			a := m.GetFieldValue(f.Field)
			b := f.Value
			ok, _ := query.EvaluateOperation(f.Op, fieldTypes[f.Field], a, b)
			if !ok {
				return false
			}
		}
		return true

	default:
		return false
	}
}

func checkOrders(q *query.Query, fieldTypes map[int]string) error {
	for _, order := range q.Orders {
		fieldType, ok := fieldTypes[order.Field]
		if !ok {
			return model.ErrQueryEntityFieldNotFound
		}
		if _, ok = query.OperatorsRegistry[fieldType]; !ok {
			return model.ErrQueryEntityFieldOperatorNotFound
		}
	}
	return nil
}

func order(mem []register.Model, orders []query.Order, fieldTypes map[int]string) {
	parsort.StructAsc(mem, func(a, b register.Model) bool {
		for _, o := range orders {
			ft := fieldTypes[o.Field]
			va := a.GetFieldValue(o.Field)
			vb := b.GetFieldValue(o.Field)

			switch o.Type {
			case query.OrderTypeAsc:
				ok, _ := query.EvaluateOperation(query.OperatorTypeLessThan, ft, va, vb)
				eq, _ := query.EvaluateOperation(query.OperatorTypeEqual, ft, va, vb)
				if !eq {
					return ok
				}
			case query.OrderTypeDesc:
				ok, _ := query.EvaluateOperation(query.OperatorTypeGreaterThan, ft, va, vb)
				eq, _ := query.EvaluateOperation(query.OperatorTypeEqual, ft, va, vb)
				if !eq {
					return ok
				}
			}
		}
		return false
	})
}
