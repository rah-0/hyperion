package template

import (
	"encoding/gob"
	"go/format"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/register"
	"github.com/rah-0/hyperion/util"

	//
	// Dynamic Imports Start
	"github.com/rah-0/hyperion/entities/Sample/v1"
	// Dynamic Imports End
	//
)

var (
	pathEntities       = filepath.Join("..", "entities")
	pathGoMod          = filepath.Join("..", "go.mod")
	entitiesRegistered = false
)

func RegisterEntities() error {
	if entitiesRegistered {
		return nil
	}

	// Types should be registered here
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
	gob.Register(register.IndexAccessor{})

	//
	// Dynamic Register Start
	if err := SampleV1.Register(); err != nil {
		return err
	}
	// Dynamic Register End
	//

	entitiesRegistered = true
	return nil
}

func Generate() error {
	pe, err := filepath.Abs(pathEntities)
	if err != nil {
		return err
	}

	structs, err := util.StructsExtractFromPackage(pe, false, 0)
	if err != nil {
		return err
	}

	err = createDirectoriesForStructs(pe, structs)
	if err != nil {
		return err
	}

	err = createEntities(structs, pe)
	if err != nil {
		return err
	}

	err = updateDynamicCode(pe)
	if err != nil {
		return err
	}

	return nil
}

func updateDynamicCode(pathEntities string) error {
	mn, err := util.GetModuleName(pathGoMod)
	if err != nil {
		return err
	}

	ds, err := util.DirectoriesInPath(pathEntities)
	if err != nil {
		return err
	}

	var imports = ""
	var registers = ""
	for i, entityName := range ds {
		versions, err := util.DirectoriesInPath(filepath.Join(pathEntities, entityName))
		if err != nil {
			return err
		}
		for _, v := range versions {
			imports += "\t" + `"` + filepath.Join(mn, "entities", entityName, v) + `"`
			registers += "\t" + "if err := " + entityName + strings.ToUpper(v) + `.Register(); err != nil { return err }`
			if i != len(ds)-1 {
				imports += "\n"
				registers += "\n"
			}
		}
	}

	err = util.FileExpand("gen.go", []util.FileExpanderTags{{
		StartTag:   []byte("// Dynamic Imports Start"),
		EndTag:     []byte("// Dynamic Imports End"),
		ExpandWith: []byte(imports),
		Count:      1,
	}})
	if err != nil {
		return err
	}

	err = util.FileExpand("gen.go", []util.FileExpanderTags{{
		StartTag:   []byte("// Dynamic Register Start"),
		EndTag:     []byte("// Dynamic Register End"),
		ExpandWith: []byte(registers),
		Count:      1,
	}})
	if err != nil {
		return err
	}

	return nil
}

func createEntityAtPath(s util.StructDef, v string, p string) error {
	var newFields []util.StructField
	newFields = append(newFields, model.GlobalStructFields...)
	newFields = append(newFields, s.Fields...)
	s.Fields = newFields

	t, err := TemplateEntity(s, v)
	if err != nil {
		return err
	}

	f, err := format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = util.FileCreate(p, f)
	if err != nil {
		return err
	}

	return nil
}

func createEntityFirstVersion(s util.StructDef, pathEntities string) error {
	p := path.Join(pathEntities, s.Name, "v1")
	err := util.DirectoryCreate(p)
	if err != nil {
		return err
	}

	p = path.Join(p, "entity.go")
	err = createEntityAtPath(s, "v1", p)
	if err != nil {
		return err
	}

	return nil
}

// isMigrationRequired checks the entity defined by the user with the latest
// version that was generated. This verifies if all fields and types match.
func isMigrationRequired(s util.StructDef, p string) (bool, string, error) {
	hv, err := util.DirectoryGetHighestVersion(p)
	if err != nil {
		return false, hv, err
	}

	p = path.Join(p, hv)
	svs, err := util.StructsExtractFromPackage(p, false, 1)
	if err != nil {
		return false, hv, err
	}

	for _, x := range svs {
		if x.Name == s.Name {
			return !compareStructFields(s, x), hv, nil
		}
	}

	return false, hv, model.ErrGeneratorStructNotFound
}

// compareStructFields compares two StructDefs while ignoring GlobalStructFields.
func compareStructFields(a, b util.StructDef) bool {
	ignoredFields := make(map[string]bool)
	for _, f := range model.GlobalStructFields {
		ignoredFields[f.Name] = true
	}

	var filteredA, filteredB []util.StructField
	for _, f := range a.Fields {
		if !ignoredFields[f.Name] {
			filteredA = append(filteredA, util.StructField{
				Name: f.Name,
				Type: f.Type,
				Tag:  "", // ignore tag
			})
		}
	}
	for _, f := range b.Fields {
		if !ignoredFields[f.Name] {
			filteredB = append(filteredB, util.StructField{
				Name: f.Name,
				Type: f.Type,
				Tag:  "", // ignore tag
			})
		}
	}

	return reflect.DeepEqual(filteredA, filteredB)
}

// createEntityMigration will create upgrade and downgrade functions together with tests.
// It is mandatory to modify the functions body and the tests
func createEntityMigration(sCurrent util.StructDef, vPrevious string, vCurrent, p string) error {
	svs, err := util.StructsExtractFromPackage(path.Join(p, vPrevious), false, 1)
	if err != nil {
		return err
	}

	var sPrevious util.StructDef
	for _, x := range svs {
		if x.Name == sCurrent.Name {
			sPrevious = x
			break
		}
	}
	if sPrevious.Name == "" {
		return model.ErrGeneratorStructNotFound
	}

	t, err := TemplateMigrations(sPrevious, sCurrent, vPrevious, vCurrent)
	if err != nil {
		return err
	}

	f, err := format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = util.FileCreate(filepath.Join(p, vCurrent, "migrations.go"), f)
	if err != nil {
		return err
	}

	t, err = TemplateMigrationsTests(sCurrent, vCurrent)
	if err != nil {
		return err
	}

	f, err = format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = util.FileCreate(filepath.Join(p, vCurrent, "migrations_test.go"), f)
	if err != nil {
		return err
	}

	return nil
}

func createEntities(structs []util.StructDef, pathEntities string) error {
	for _, s := range structs {
		p := path.Join(pathEntities, s.Name)
		versions, err := util.DirectoriesInPath(p)
		if err != nil {
			return err
		}

		if len(versions) == 0 {
			err = createEntityFirstVersion(s, pathEntities)
			if err != nil {
				return err
			}
		} else {
			mr, hv, err := isMigrationRequired(s, p)
			if err != nil {
				return err
			}

			if mr {
				nv, err := util.StringNextVersion(hv)
				if err != nil {
					return err
				}

				err = createEntityAtPath(s, nv, path.Join(p, nv, "entity.go"))
				if err != nil {
					return err
				}

				err = createEntityMigration(s, hv, nv, p)
				if err != nil {
					return err
				}
			} else {
				err = createEntityAtPath(s, hv, path.Join(p, hv, "entity.go"))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createDirectoriesForStructs(pathEntities string, structs []util.StructDef) (err error) {
	for _, s := range structs {
		pathEntity := path.Join(pathEntities, s.Name)
		err = util.DirectoryCreate(pathEntity)
		if err != nil {
			return
		}
	}
	return
}
