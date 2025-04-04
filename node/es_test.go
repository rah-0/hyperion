package node

import (
	"testing"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
)

func TestOrderMultiField(t *testing.T) {
	a := &SampleV1.Sample{Name: "B"}
	b := &SampleV1.Sample{Name: "A"}

	mem := []register.Model{a, b}

	order(mem, []query.Order{
		{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
	}, SampleV1.FieldTypes)

	expected := []string{"A", "B"}
	for i, m := range mem {
		got := m.(*SampleV1.Sample).Name
		if got != expected[i] {
			t.Fatalf("at %d: got %s, want %s", i, got, expected[i])
		}
	}
}

func TestOrderSortingCombinations(t *testing.T) {
	cases := []struct {
		name     string
		entities []*SampleV1.Sample
		orders   []query.Order
		expected []string
	}{
		{
			name: "Asc on single field",
			entities: []*SampleV1.Sample{
				{Name: "B"}, {Name: "A"},
			},
			orders: []query.Order{
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
			},
			expected: []string{"A", "B"},
		},
		{
			name: "Desc on single field",
			entities: []*SampleV1.Sample{
				{Name: "B"}, {Name: "A"},
			},
			orders: []query.Order{
				{Type: query.OrderTypeDesc, Field: SampleV1.FieldName},
			},
			expected: []string{"B", "A"},
		},
		{
			name: "Asc on multiple fields",
			entities: []*SampleV1.Sample{
				{Name: "A", Surname: "Z"},
				{Name: "A", Surname: "B"},
				{Name: "B", Surname: "A"},
			},
			orders: []query.Order{
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldSurname},
			},
			expected: []string{"A|B", "A|Z", "B|A"},
		},
		{
			name: "Desc on multiple fields",
			entities: []*SampleV1.Sample{
				{Name: "A", Surname: "Z"},
				{Name: "A", Surname: "B"},
				{Name: "B", Surname: "A"},
			},
			orders: []query.Order{
				{Type: query.OrderTypeDesc, Field: SampleV1.FieldName},
				{Type: query.OrderTypeDesc, Field: SampleV1.FieldSurname},
			},
			expected: []string{"B|A", "A|Z", "A|B"},
		},
		{
			name: "Mixed Asc/Desc",
			entities: []*SampleV1.Sample{
				{Name: "A", Surname: "Z"},
				{Name: "A", Surname: "B"},
				{Name: "B", Surname: "A"},
			},
			orders: []query.Order{
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
				{Type: query.OrderTypeDesc, Field: SampleV1.FieldSurname},
			},
			expected: []string{"A|Z", "A|B", "B|A"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mem := make([]register.Model, len(c.entities))
			for i, e := range c.entities {
				mem[i] = e
			}

			order(mem, c.orders, SampleV1.FieldTypes)

			for i, m := range mem {
				got := m.(*SampleV1.Sample)
				expect := c.expected[i]
				actual := got.Name
				if got.Surname != "" {
					actual += "|" + got.Surname
				}
				if actual != expect {
					t.Fatalf("at %d: got %s, want %s", i, actual, expect)
				}
			}
		})
	}
}
