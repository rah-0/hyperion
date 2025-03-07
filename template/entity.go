package template

import (
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/rah-0/hyperion/util"
)

func TemplateEntity(s StructDef, v string) (string, error) {
	mn, err := GetModuleName("go.mod")
	if err != nil {
		return "", err
	}

	template := "package " + s.Name + "\n"
	template += "// The code in this file is autogenerated, do not modify manually!" + "\n\n"

	template += "import (\n"
	template += `"encoding/gob"` + "\n\n"
	template += `. "` + filepath.Join(mn, "register") + `"` + "\n"
	template += ")\n"

	template += "const (\n"
	template += `Version = "` + v + `"` + "\n"
	template += `EntityName = "` + s.Name + `"` + "\n"
	template += `DbFileName = "` + s.Name + strings.ToUpper(v) + `.bin"` + "\n"
	template += ")\n\n"

	i := 1
	template += "var Fields = map[string]int{" + "\n"
	for _, f := range s.Fields {
		template += `"` + f.Name + `": ` + strconv.Itoa(i) + "," + "\n"
		i++
	}
	template += "}" + "\n\n"

	template += "var _ Model = (*" + s.Name + ")(nil)" + "\n"

	template += "func init(){\n"
	template += "gob.Register(" + s.Name + "{})\n\n"
	template += "RegisterEntity(&Entity{" + "\n"
	template += "Version: Version," + "\n"
	template += "EntityName: EntityName," + "\n"
	template += "DbFileName: DbFileName," + "\n"
	template += "Fields: Fields," + "\n"
	template += "New: New," + "\n"
	template += "})" + "\n"
	template += "}\n"

	template += "type " + s.Name + " struct {\n"
	for _, f := range s.Fields {
		template += f.Name + " " + f.Type + "\n"
	}
	template += "}\n"

	template += "func New() Model {" + "\n"
	template += "return &" + s.Name + "{}\n"
	template += "}\n\n"

	template += "func (s *" + s.Name + ") SetFieldValue(fieldName string, value any) {" + "\n"
	template += "switch Fields[fieldName] {" + "\n"
	i = 1
	for _, f := range s.Fields {
		template += "case " + strconv.Itoa(i) + ":\n"
		template += "if v, ok := value.(" + f.Type + "); ok {" + "\n"
		template += "s." + f.Name + "=v" + "\n"
		template += "}" + "\n"
		i++
	}
	template += "}}" + "\n\n"

	template += "func (s *" + s.Name + ") GetFieldValue(fieldName string) any {" + "\n"
	template += "switch Fields[fieldName] {" + "\n"
	i = 1
	for _, f := range s.Fields {
		template += "case " + strconv.Itoa(i) + ":\n"
		template += "return s." + f.Name + "\n"
		i++
	}
	template += "}" + "\n"
	template += "return nil" + "\n"
	template += "}" + "\n"

	return template, nil
}

func TemplateMigrations(sPrevious StructDef, sCurrent StructDef, vPrevious string, vCurrent string) (string, error) {
	mn, err := GetModuleName("go.mod")
	if err != nil {
		return "", err
	}

	template := "package " + sCurrent.Name + "\n"
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

func TemplateMigrationsTests(sCurrent StructDef) (string, error) {
	template := "package " + sCurrent.Name + "\n"
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
