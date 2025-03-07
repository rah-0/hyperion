package register

var Entities []*Entity

type Entity struct {
	Version    string
	EntityName string
	DbFileName string
	Fields     map[string]int
	New        func() Model
}

func RegisterEntity(entity *Entity) {
	Entities = append(Entities, entity)
}

type Model interface {
	SetFieldValue(fieldName string, value any)
	GetFieldValue(fieldName string) any
}
