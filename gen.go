package main

import (
	"go/format"
	"path"
	"path/filepath"
	"reflect"

	. "github.com/rah-0/hyperion/template"
	. "github.com/rah-0/hyperion/util"
	//
	// Dynamic Imports Start
	_ "github.com/rah-0/hyperion/entities/Sample/v1"
	// Dynamic Imports End
)

func Generate() error {
	pathEntities, err := filepath.Abs("entities")
	if err != nil {
		return err
	}

	structs, err := StructsExtractFromPackage(pathEntities, false)
	if err != nil {
		return err
	}

	err = createDirectoriesForStructs(pathEntities, structs)
	if err != nil {
		return err
	}

	err = createEntities(structs, pathEntities)
	if err != nil {
		return err
	}

	err = updateDynamicImports(pathEntities)
	if err != nil {
		return err
	}

	return nil
}

func updateDynamicImports(pathEntities string) error {
	mn, err := GetModuleName("go.mod")
	if err != nil {
		return err
	}

	ds, err := DirectoriesInPath(pathEntities)
	if err != nil {
		return err
	}

	var imports = ""
	for _, entityName := range ds {
		versions, err := DirectoriesInPath(filepath.Join(pathEntities, entityName))
		if err != nil {
			return err
		}
		for i, v := range versions {
			imports += "\t" + `_ "` + filepath.Join(mn, "entities", entityName, v) + `"`
			if i < len(versions)-1 {
				imports += "\n"
			}
		}
	}

	err = FileExpand("gen.go", []FileExpanderTags{{
		StartTag:   []byte("// Dynamic Imports Start"),
		EndTag:     []byte("// Dynamic Imports End"),
		ExpandWith: []byte(imports),
		Count:      1,
	}})
	if err != nil {
		return err
	}

	return nil
}

func createEntityAtPath(s StructDef, v string, p string) error {
	t, err := TemplateEntity(s, v)
	if err != nil {
		return err
	}

	f, err := format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = FileCreate(p, f)
	if err != nil {
		return err
	}

	return nil
}

func createEntityFirstVersion(s StructDef, pathEntities string) error {
	p := path.Join(pathEntities, s.Name, "v1")
	err := DirectoryCreate(p)
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
func isMigrationRequired(s StructDef, p string) (bool, string, error) {
	hv, err := DirectoryGetHighestVersion(p)
	if err != nil {
		return false, hv, err
	}

	p = path.Join(p, hv)
	svs, err := StructsExtractFromPackage(p, false)
	if err != nil {
		return false, hv, err
	}

	for _, x := range svs {
		if x.Name == s.Name {
			return !reflect.DeepEqual(s, x), hv, nil
		}
	}

	return false, hv, ErrGeneratorStructNotFound
}

// createEntityMigration will create upgrade and downgrade functions together with tests.
// It is mandatory to modify the functions body and the tests
func createEntityMigration(sCurrent StructDef, vPrevious string, vCurrent, p string) error {
	svs, err := StructsExtractFromPackage(path.Join(p, vPrevious), false)
	if err != nil {
		return err
	}

	var sPrevious StructDef
	for _, x := range svs {
		if x.Name == sCurrent.Name {
			sPrevious = x
			break
		}
	}
	if sPrevious.Name == "" {
		return ErrGeneratorStructNotFound
	}

	t, err := TemplateMigrations(sPrevious, sCurrent, vPrevious, vCurrent)
	if err != nil {
		return err
	}

	f, err := format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = FileCreate(filepath.Join(p, vCurrent, "migrations.go"), f)
	if err != nil {
		return err
	}

	t, err = TemplateMigrationsTests(sCurrent)
	if err != nil {
		return err
	}

	f, err = format.Source([]byte(t))
	if err != nil {
		return err
	}

	err = FileCreate(filepath.Join(p, vCurrent, "migrations_test.go"), f)
	if err != nil {
		return err
	}

	return nil
}

func createEntities(structs []StructDef, pathEntities string) error {
	for _, s := range structs {
		p := path.Join(pathEntities, s.Name)
		versions, err := DirectoriesInPath(p)
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
				nv, err := StringNextVersion(hv)
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

func createDirectoriesForStructs(pathEntities string, structs []StructDef) (err error) {
	for _, s := range structs {
		pathEntity := path.Join(pathEntities, s.Name)
		err = DirectoryCreate(pathEntity)
		if err != nil {
			return
		}
	}
	return
}
