package node

import (
	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
)

func (x *EntityStorage) HandleQuery(q *query.Query) ([]register.Model, error) {
	if q == nil {
		return nil, model.ErrQueryNil
	}

	var results []register.Model
	mem := x.Memory.New().MemoryGetAll()
	filters := q.Filters.Filters
	filterType := q.Filters.Type

	for _, m := range mem {
		match := false

		switch filterType {
		case query.FilterTypeOr:
			for _, f := range filters {
				fieldType := x.Entity.FieldTypes[f.Field]
				a := m.GetFieldValue(f.Field)
				b := f.Value
				ok, err := query.EvaluateOp(f.Op, fieldType, a, b)
				if err != nil {
					return nil, err
				}
				if ok {
					match = true
					break
				}
			}

		case query.FilterTypeAnd:
			match = true
			for _, f := range filters {
				fieldType := x.Entity.FieldTypes[f.Field]
				a := m.GetFieldValue(f.Field)
				b := f.Value
				ok, err := query.EvaluateOp(f.Op, fieldType, a, b)
				if err != nil {
					return nil, err
				}
				if !ok {
					match = false
					break
				}
			}

		default:
			return nil, model.ErrQueryInvalidFilterType
		}

		if match {
			results = append(results, m)
		}
	}

	return results, nil
}
