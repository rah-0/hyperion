package template

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rah-0/hyperion/util"
)

func TemplateEntity(s util.StructDef, v string) (string, error) {
	mn, err := util.GetModuleName(pathGoMod)
	if err != nil {
		return "", err
	}

	template := "package " + s.Name + strings.ToUpper(v) + "\n"
	template += "// ---------------------------------------------------------------" + "\n"
	template += "// The code in this file is autogenerated, do not modify manually!" + "\n"
	template += "// ---------------------------------------------------------------" + "\n\n"

	template += "import (\n"
	template += `"bytes"` + "\n"
	template += `"encoding/gob"` + "\n"
	template += `"errors"` + "\n"
	template += `"sync"` + "\n"
	for _, f := range s.Fields {
		if f.Type == "time.Time" {
			template += `"time"` + "\n"
		}
	}
	template += "\n"
	template += `"github.com/google/uuid"` + "\n\n"
	template += `"` + filepath.Join(mn, "hconn") + `"` + "\n"
	template += `"` + filepath.Join(mn, "model") + `"` + "\n"
	template += `"` + filepath.Join(mn, "query") + `"` + "\n"
	template += `"` + filepath.Join(mn, "register") + `"` + "\n"
	template += ")\n\n"

	template += "const (\n"
	template += `Version = "` + v + `"` + "\n"
	template += `Name = "` + s.Name + `"` + "\n"
	template += `DbFileName = "` + s.Name + strings.ToUpper(v) + `.bin"` + "\n"
	template += ")\n\n"

	i := 1
	template += "const (\n"
	for _, f := range s.Fields {
		template += `Field` + f.Name + ` = ` + strconv.Itoa(i) + "\n"
		i++
	}
	template += ")\n\n"

	i = 1
	template += "var FieldTypes = map[int]string{" + "\n"
	for _, f := range s.Fields {
		template += `Field` + f.Name + `: "` + f.Type + `",` + "\n"
		i++
	}
	template += "}" + "\n\n"

	template += "var Indexes = map[int]any{\n"
	for _, f := range s.Fields {
		template += "\tField" + f.Name + ": map[" + f.Type + "][]*" + s.Name + "{},\n"
	}
	template += "}\n\n"

	template += "var (" + "\n"
	template += "_ register.Model = (*" + s.Name + ")(nil)" + "\n"
	template += "mu sync.Mutex" + "\n"
	template += "Buffer = new(bytes.Buffer)" + "\n"
	template += "Encoder = gob.NewEncoder(Buffer)" + "\n"
	template += "Decoder = gob.NewDecoder(Buffer)" + "\n"
	template += "Mem []*" + s.Name + "\n"
	template += "IndexAccessors = map[int]register.IndexAccessor{}\n"
	template += ")" + "\n\n"

	template += "func init() {\n"
	template += "// Validate all FieldTypes have an operator set\n"
	template += "for _, typ := range FieldTypes {\n"
	template += "if _, ok := query.OperatorsRegistry[typ]; !ok {\n"
	template += `panic("missing operator set for field type: " + typ)` + "\n"
	template += "}\n"
	template += "}\n\n"
	template += "//The following process initializes the encoder and decoder by preloading metadata." + "\n"
	template += "//This prevents metadata from being stored with the first encoded struct." + "\n"
	template += "//If the metadata were missing or inconsistent, decoding the struct later could fail." + "\n"
	template += "gob.Register(&" + s.Name + "{})\n"
	template += "x := New()\n"
	template += "if err := x.Encode(); err != nil {\n"
	template += `panic("failed to encode type metadata: " + err.Error())` + "\n"
	template += "}\n"
	template += "if err := x.Decode(); err != nil {\n"
	template += `panic("failed to decode type metadata: " + err.Error())` + "\n"
	template += "}\n"
	template += "x.BufferReset()\n\n"
	template += "// IndexAccessors definitions" + "\n"
	for _, f := range s.Fields {
		template += "IndexAccessors[Field" + f.Name + "] = func(val any) []register.Model {\n"
		template += "idx := Indexes[Field" + f.Name + "].(map[" + f.Type + "][]*" + s.Name + ")\n"
		template += "v, ok := val.(" + f.Type + ")\n"
		template += "if !ok {\n"
		template += "return nil\n"
		template += "}\n"
		template += "return CastToModel(idx[v])\n"
		template += "}\n"
	}
	template += "\n"

	template += "// Initializations" + "\n"
	template += "Mem = []*" + s.Name + "{}\n"
	template += "register.RegisterEntity(\n"
	template += "&register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "DbFileName: DbFileName,\n"
	template += "}, &register.EntityExtension{\n"
	template += "New: New,\n"
	template += "FieldTypes: FieldTypes,\n"
	template += "Indexes: Indexes,\n"
	template += "IndexAccessors: IndexAccessors,\n"
	template += "},\n"
	template += ")\n"
	template += "}\n"

	template += "type " + s.Name + " struct {\n"
	for _, f := range s.Fields {
		template += f.Name + " " + f.Type + " " + f.Tag + "\n"
	}
	template += "}\n"

	template += "func New() register.Model {\n"
	template += "return &" + s.Name + "{}\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") WithNewUuid() {\n"
	template += "s.Uuid = uuid.New()\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") SetUuid(uuid uuid.UUID) {\n"
	template += "s.Uuid = uuid\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") GetUuid() uuid.UUID {\n"
	template += "return s.Uuid\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") IsDeleted() bool {\n"
	template += "return s.Deleted\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") SetFieldValue(field int, value any) {" + "\n"
	template += "switch field {" + "\n"
	i = 1
	for _, f := range s.Fields {
		template += "case Field" + f.Name + ":\n"
		template += "if v, ok := value.(" + f.Type + "); ok {" + "\n"
		template += "s." + f.Name + "=v" + "\n"
		template += "}" + "\n"
		i++
	}
	template += "}}" + "\n\n"

	template += "func (s *" + s.Name + ") GetFieldValue(field int) any {" + "\n"
	template += "switch field {" + "\n"
	i = 1
	for _, f := range s.Fields {
		template += "case Field" + f.Name + ":\n"
		template += "return s." + f.Name + "\n"
		i++
	}
	template += "}" + "\n"
	template += "return nil" + "\n"
	template += "}" + "\n\n"

	template += "func (s *" + s.Name + ") Encode() error {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "return Encoder.Encode(s)\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") Decode() error {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "return Decoder.Decode(s)\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") BufferReset() {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "Buffer.Reset()\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") GetBuffer() *bytes.Buffer {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "return Buffer\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") GetBufferData() []byte {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "return Buffer.Bytes()\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") SetBufferData(data []byte) {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "Buffer.Write(data)\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryAdd() {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "Mem = append(Mem, s)\n\n"
	template += "// Update indexes\n"
	for _, f := range s.Fields {
		varName := "index" + f.Name
		template += varName + " := Indexes[Field" + f.Name + "].(map[" + f.Type + "][]*" + s.Name + ")\n"
		template += varName + "[s." + f.Name + "] = append(" + varName + "[s." + f.Name + "], s)\n"
	}
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryRemove() {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "for i, instance := range Mem {\n"
	template += "if instance == s {\n"
	template += "lastIndex := len(Mem) - 1\n"
	template += "Mem[i] = Mem[lastIndex]\n"
	template += "Mem = Mem[:lastIndex]\n"
	template += "break\n"
	template += "}\n"
	template += "}\n\n"
	template += "// Remove from indexes\n"
	for _, f := range s.Fields {
		varName := "index" + f.Name
		template += varName + " := Indexes[Field" + f.Name + "].(map[" + f.Type + "][]*" + s.Name + ")\n"
		template += varName + "[s." + f.Name + "] = removeFromIndex(" + varName + "[s." + f.Name + "], s)\n"
	}
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryUpdate() {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n\n"
	template += "for i, old := range Mem {\n"
	template += "if old.Uuid != s.Uuid {\n"
	template += "continue\n"
	template += "}\n\n"
	template += "// Update indexes only if values changed\n"
	for _, f := range s.Fields {
		template += "if old." + f.Name + " != s." + f.Name + " {\n"
		varName := "index" + f.Name
		template += varName + " := Indexes[Field" + f.Name + "].(map[" + f.Type + "][]*" + s.Name + ")\n"
		template += varName + "[old." + f.Name + "] = removeFromIndex(" + varName + "[old." + f.Name + "], old)\n"
		template += varName + "[s." + f.Name + "] = append(" + varName + "[s." + f.Name + "], s)\n"
		template += "}\n"
	}
	template += "\nMem[i] = s\n"
	template += "break\n"
	template += "}\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryClear() {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n\n"
	template += "Mem = []*" + s.Name + "{}\n\n"
	for _, f := range s.Fields {
		template += "Indexes[Field" + f.Name + "] = map[" + f.Type + "][]*" + s.Name + "{}\n"
	}
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryGetAll() []register.Model {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n"
	template += "instances := make([]register.Model, len(Mem))\n"
	template += "for i, instance := range Mem {\n"
	template += "instances[i] = instance\n"
	template += "}\n"
	template += "return instances\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") MemoryContains(target register.Model) bool {\n"
	template += "mu.Lock()\n"
	template += "defer mu.Unlock()\n\n"
	template += "for _, instance := range Mem {\n"
	template += "if instance == target {\n"
	template += "return true\n"
	template += "}\n"
	template += "}\n"
	template += "return false\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") DbInsert(c *hconn.HConn) error {\n"
	template += "if s.Uuid == uuid.Nil {\n"
	template += "s.WithNewUuid()\n"
	template += "}\n"
	template += "if err := s.Encode(); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "msg := model.Message{\n"
	template += "Type: model.MessageTypeInsert,\n"
	template += "Entity: register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "Data: s.GetBufferData(),\n"
	template += "},\n"
	template += "}\n"
	template += "s.BufferReset()\n\n"
	template += "if err := c.Send(msg); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "resp, err := c.Receive()\n"
	template += "if err != nil {\n"
	template += "return err\n"
	template += "}\n"
	template += "if resp.Status == model.StatusError {\n"
	template += "return errors.New(resp.String)\n"
	template += "}\n\n"
	template += "return nil\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") DbDelete(c *hconn.HConn) error {\n"
	template += "s.Deleted = true\n"
	template += "if err := s.Encode(); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "msg := model.Message{\n"
	template += "Type: model.MessageTypeDelete,\n"
	template += "Entity: register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "Data: s.GetBufferData(),\n"
	template += "},\n"
	template += "}\n"
	template += "s.BufferReset()\n\n"
	template += "if err := c.Send(msg); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "resp, err := c.Receive()\n"
	template += "if err != nil {\n"
	template += "return err\n"
	template += "}\n"
	template += "if resp.Status == model.StatusError {\n"
	template += "return errors.New(resp.String)\n"
	template += "}\n\n"
	template += "return nil\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") DbUpdate(c *hconn.HConn) error {\n"
	template += "if s.Uuid == uuid.Nil {\n"
	template += `return model.ErrQueryEntityNoUuid` + "\n"
	template += "}\n"
	template += "if err := s.Encode(); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "msg := model.Message{\n"
	template += "Type: model.MessageTypeUpdate,\n"
	template += "Entity: register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "Data: s.GetBufferData(),\n"
	template += "},\n"
	template += "}\n"
	template += "s.BufferReset()\n\n"
	template += "if err := c.Send(msg); err != nil {\n"
	template += "return err\n"
	template += "}\n\n"
	template += "resp, err := c.Receive()\n"
	template += "if err != nil {\n"
	template += "return err\n"
	template += "}\n"
	template += "if resp.Status == model.StatusError {\n"
	template += "return errors.New(resp.String)\n"
	template += "}\n\n"
	template += "return nil\n"
	template += "}\n\n"

	template += "func DbGetAll(c *hconn.HConn) ([]*" + s.Name + ", error) {\n"
	template += "msg := model.Message{\n"
	template += "Type: model.MessageTypeGetAll,\n"
	template += "Entity: register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "},\n"
	template += "}\n\n"
	template += "if err := c.Send(msg); err != nil {\n"
	template += "return nil, err\n"
	template += "}\n\n"
	template += "resp, err := c.Receive()\n"
	template += "if err != nil {\n"
	template += "return nil, err\n"
	template += "}\n\n"
	template += "if resp.Status == model.StatusError {\n"
	template += "return nil, errors.New(resp.String)\n"
	template += "}\n\n"
	template += "return CastTo" + s.Name + "(resp.Models), nil\n"
	template += "}\n\n"

	template += "func DbQuery(c *hconn.HConn, q *query.Query) ([]*" + s.Name + ", error) {\n"
	template += "msg := model.Message{\n"
	template += "Type: model.MessageTypeQuery,\n"
	template += "Entity: register.EntityBase{\n"
	template += "Version: Version,\n"
	template += "Name: Name,\n"
	template += "},\n"
	template += "Query: q,\n"
	template += "}\n\n"
	template += "if err := c.Send(msg); err != nil {\n"
	template += "return nil, err\n"
	template += "}\n\n"
	template += "resp, err := c.Receive()\n"
	template += "if err != nil {\n"
	template += "return nil, err\n"
	template += "}\n\n"
	template += "if resp.Status == model.StatusError {\n"
	template += "return nil, errors.New(resp.String)\n"
	template += "}\n\n"
	template += "return CastTo" + s.Name + "(resp.Models), nil\n"
	template += "}\n\n"

	template += "func CastTo" + s.Name + "(models []register.Model) []*" + s.Name + " {\n"
	template += "out := make([]*" + s.Name + ", len(models))\n"
	template += "for i, m := range models {\n"
	template += "out[i] = m.(*" + s.Name + ")\n"
	template += "}\n"
	template += "return out\n"
	template += "}\n\n"

	template += "func CastToModel(items []*" + s.Name + ") []register.Model {\n"
	template += "out := make([]register.Model, len(items))\n"
	template += "for i, s := range items {\n"
	template += "out[i] = s\n"
	template += "}\n"
	template += "return out\n"
	template += "}\n\n"

	template += "func removeFromIndex(list []*" + s.Name + ", target *" + s.Name + ") []*" + s.Name + " {\n"
	template += "for i, item := range list {\n"
	template += "if item == target {\n"
	template += "last := len(list) - 1\n"
	template += "list[i] = list[last]\n"
	template += "return list[:last]\n"
	template += "}\n"
	template += "}\n"
	template += "return list\n"
	template += "}\n\n"

	return template, nil
}

func TemplateMigrations(sPrevious util.StructDef, sCurrent util.StructDef, vPrevious string, vCurrent string) (string, error) {
	mn, err := util.GetModuleName(pathGoMod)
	if err != nil {
		return "", err
	}

	template := "package " + sCurrent.Name + strings.ToUpper(vCurrent) + "\n"
	template += "// NOTE: this file is generated only once, if you want to update it you can delete it and run the generator again!" + "\n\n"

	template += "import (\n"
	template += vPrevious + ` "` + filepath.Join(mn, "entities", sCurrent.Name, vPrevious) + `"`
	template += ")\n"

	template += `func Upgrade(previous ` + vPrevious + `.` + sPrevious.Name + `) (current ` + sCurrent.Name + `){` + "\n"
	template += `panic("Function not implemented")` + "\n"
	template += "}\n\n"

	template += `func Downgrade(current ` + sCurrent.Name + `) (previous ` + vPrevious + `.` + sPrevious.Name + `){` + "\n"
	template += `panic("Function not implemented")` + "\n"
	template += "}\n"

	return template, nil
}

func TemplateMigrationsTests(sCurrent util.StructDef, vCurrent string) (string, error) {
	template := "package " + sCurrent.Name + strings.ToUpper(vCurrent) + "\n"
	template += "// NOTE: this file is generated only once, if you want to update it you can delete it and run the generator again!" + "\n\n"

	template += "import (\n"
	template += `"testing"`
	template += ")\n"

	template += `func TestUpgrade(t *testing.T) {` + "\n"
	template += `t.Fatal("Test not implemented")` + "\n"
	template += "}\n\n"

	template += `func TestDowngrade(t *testing.T) {` + "\n"
	template += `t.Fatal("Test not implemented")` + "\n"
	template += "}\n\n"

	return template, nil
}
