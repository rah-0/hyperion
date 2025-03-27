package node

import (
	"runtime"
	"sync"

	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
)

func (x *EntityStorage) HandleQuery(q *query.Query) ([]register.Model, error) {
	mem := x.Memory.New().MemoryGetAll()
	total := len(mem)
	if total == 0 {
		return nil, model.ErrQueryMemoryEmpty
	}
	if q == nil {
		return nil, model.ErrQueryNil
	}
	if q.Filters.Type == query.FilterTypeUndefined {
		return nil, model.ErrQueryInvalidFilterType
	}

	fieldTypes := x.Memory.FieldTypes
	filters := q.Filters.Filters
	filterType := q.Filters.Type
	limit := q.Limit

	numWorkers := runtime.NumCPU()
	if total < numWorkers {
		numWorkers = total
	}

	chunkSize := total / numWorkers
	remainder := total % numWorkers

	var wg sync.WaitGroup
	resultsCh := make(chan register.Model)

	start := 0
	for i := 0; i < numWorkers; i++ {
		end := start + chunkSize
		if i < remainder {
			end++
		}
		part := mem[start:end]
		start = end

		wg.Add(1)
		go func(part []register.Model) {
			defer wg.Done()
			for _, m := range part {
				if matchModel(m, filters, fieldTypes, filterType) {
					resultsCh <- m
				}
			}
		}(part)
	}
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var results []register.Model
	for m := range resultsCh {
		results = append(results, m)
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, nil
}

func matchModel(m register.Model, filters []query.Filter, fieldTypes map[int]string, typ query.FilterType) bool {
	switch typ {
	case query.FilterTypeOr:
		for _, f := range filters {
			a := m.GetFieldValue(f.Field)
			b := f.Value
			ok, _ := query.EvaluateOp(f.Op, fieldTypes[f.Field], a, b)
			if ok {
				return true
			}
		}
		return false

	case query.FilterTypeAnd:
		for _, f := range filters {
			a := m.GetFieldValue(f.Field)
			b := f.Value
			ok, _ := query.EvaluateOp(f.Op, fieldTypes[f.Field], a, b)
			if !ok {
				return false
			}
		}
		return true

	default:
		return false
	}
}
