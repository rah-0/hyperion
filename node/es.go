package node

import (
	"github.com/rah-0/parsort"

	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
)

func (x *EntityStorage) HandleQuery(q *query.Query) ([]register.Model, error) {
	if q == nil {
		return nil, model.ErrQueryNil
	}

	fieldTypes := x.Memory.EntityExtension.FieldTypes
	indexAccessors := x.Memory.EntityExtension.IndexAccessors
	indexesSorted := map[int]bool{}
	for _, f := range x.Memory.EntityExtension.IndexesSorted {
		indexesSorted[f] = true
	}

	hasFilters := q.Filters.Type != query.FilterTypeUndefined
	filters := q.Filters.Filters
	filterType := q.Filters.Type

	hasLimit := q.Limit > 0
	limit := q.Limit

	hasOrders := len(q.Orders) > 0

	// Optimization: use sorted index when all ORDER BY fields are sorted
	if !hasFilters && hasOrders {
		allSorted := true
		for _, o := range q.Orders {
			if !indexesSorted[o.Field] {
				allSorted = false
				break
			}
		}
		if allSorted {
			mem := x.Memory.EntityExtension.New().MemoryGetAll()
			order(mem, q.Orders, fieldTypes)
			if hasLimit && len(mem) > limit {
				mem = mem[:limit]
			}
			return mem, nil
		}
	}

	var results []register.Model
	if hasFilters && filterType == query.FilterTypeAnd {
		bestIndex := getBestIndex(filters, indexAccessors)
		if bestIndex != nil {
			for _, m := range bestIndex {
				if matchModel(m, filters, fieldTypes, filterType) {
					results = append(results, m)
				}
			}
			if hasOrders {
				err := checkOrders(q, fieldTypes)
				if err != nil {
					return nil, err
				}
				order(results, q.Orders, fieldTypes)
			}
			if hasLimit && len(results) > limit {
				results = results[:limit]
			}
			return results, nil
		}
	}

	mem := x.Memory.EntityExtension.New().MemoryGetAll()
	for _, m := range mem {
		if !hasFilters || matchModel(m, filters, fieldTypes, filterType) {
			results = append(results, m)
		}
	}

	if hasOrders {
		err := checkOrders(q, fieldTypes)
		if err != nil {
			return nil, err
		}
		order(results, q.Orders, fieldTypes)
	}

	if hasLimit && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func getBestIndex(filters []query.Filter, indexAccessors map[int]register.IndexAccessor) []register.Model {
	var best []register.Model
	for _, f := range filters {
		if f.Op != query.OperatorTypeEqual {
			continue
		}
		get := indexAccessors[f.Field]
		if get == nil {
			continue
		}
		candidates := get(f.Value)
		if candidates == nil {
			continue
		}
		if best == nil || len(candidates) < len(best) {
			best = candidates
		}
	}
	return best
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

func applyOrdering(results []register.Model, q *query.Query, fieldTypes map[int]string) ([]register.Model, error) {
	if len(q.Orders) == 0 {
		return results, nil
	}
	if err := checkOrders(q, fieldTypes); err != nil {
		return results, err
	}
	order(results, q.Orders, fieldTypes)
	return results, nil
}

func checkOrders(q *query.Query, fieldTypes map[int]string) error {
	for _, o := range q.Orders {
		fieldType, ok := fieldTypes[o.Field]
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
				ok, _ := query.EvaluateOperation(query.OperatorTypeLesser, ft, va, vb)
				eq, _ := query.EvaluateOperation(query.OperatorTypeEqual, ft, va, vb)
				if !eq {
					return ok
				}
			case query.OrderTypeDesc:
				ok, _ := query.EvaluateOperation(query.OperatorTypeGreater, ft, va, vb)
				eq, _ := query.EvaluateOperation(query.OperatorTypeEqual, ft, va, vb)
				if !eq {
					return ok
				}
			}
		}
		return false
	})
}
