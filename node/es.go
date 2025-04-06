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

	var results []register.Model
	fieldTypes := x.Memory.EntityExtension.FieldTypes
	indexAccessors := x.Memory.EntityExtension.IndexAccessors

	hasFilters := q.Filters.Type != query.FilterTypeUndefined
	filters := q.Filters.Filters
	filterType := q.Filters.Type

	hasLimit := q.Limit > 0
	limit := q.Limit

	hasOrders := len(q.Orders) > 0
	if hasOrders {
		err := checkOrders(q, fieldTypes)
		if err != nil {
			return nil, err
		}
	}

	if hasFilters {
		var sets [][]register.Model
		if filterType == query.FilterTypeAnd {
			for _, f := range filters {
				if f.Op == query.OperatorTypeEqual {
					sets = append(sets, indexAccessors[f.Field].GetByValue(f.Value))
				} else {
					set, err := indexAccessors[f.Field].Evaluate(f.Op, f.Value)
					if err != nil {
						return nil, err
					}
					sets = append(sets, set)
				}
			}
			results = intersectSets(sets)
		} else if filterType == query.FilterTypeOr {
			for _, f := range filters {
				var set []register.Model
				if f.Op == query.OperatorTypeEqual {
					set = indexAccessors[f.Field].GetByValue(f.Value)
				} else {
					var err error
					set, err = indexAccessors[f.Field].Evaluate(f.Op, f.Value)
					if err != nil {
						return nil, err
					}
				}
				sets = append(sets, set)
			}
			results = unionSets(sets)
		}
	} else {
		results = x.Memory.EntityExtension.New().MemoryGetAll()
	}

	if hasOrders {
		parsort.StructAsc(results, func(a, b register.Model) bool {
			for _, o := range q.Orders {
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

	if hasLimit && len(results) > limit {
		results = results[:limit]
	}

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

func intersectSets(sets [][]register.Model) []register.Model {
	if len(sets) == 0 {
		return []register.Model{}
	}
	if len(sets) == 1 {
		return sets[0]
	}

	// Build map from first set
	ref := make(map[register.Model]int)
	for _, m := range sets[0] {
		ref[m] = 1
	}

	// Count occurrences in other sets
	for i := 1; i < len(sets); i++ {
		for _, m := range sets[i] {
			if ref[m] == i {
				ref[m]++
			}
		}
	}

	// Collect models that appeared in all sets
	var out []register.Model
	for m, count := range ref {
		if count == len(sets) {
			out = append(out, m)
		}
	}
	return out
}

func unionSets(sets [][]register.Model) []register.Model {
	seen := make(map[register.Model]struct{})
	var out []register.Model
	for _, set := range sets {
		for _, m := range set {
			if _, ok := seen[m]; !ok {
				seen[m] = struct{}{}
				out = append(out, m)
			}
		}
	}
	return out
}
