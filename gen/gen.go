package gen

import (
	"go/format"
	"path"
	"path/filepath"
	"reflect"

	. "github.com/rah-0/hyperion/util"
)

func Generate() error {
	pathEntities, err := filepath.Abs(filepath.Join("..", "entities"))
	if err != nil {
		return err
	}

	existingEntitiesDirs, err := DirectoriesInPath(pathEntities)
	if err != nil {
		return err
	}

	structs, err := StructsExtractFromPackage(pathEntities)
	if err != nil {
		return err
	}

	pathsEntities, err := createDirectoriesForStructs(pathEntities, structs)
	if err != nil {
		return err
	}

	err = createEntities(structs, pathEntities)
	if err != nil {
		return err
	}

	// Cleanup
	err = deleteDirectoriesOfNonExistingEntities(pathEntities, existingEntitiesDirs, pathsEntities)
	if err != nil {
		return err
	}

	return nil
}

func createEntityAtPath(s StructDef, v string, p string) error {
	t := templateEntity(s, v)
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
	svs, err := StructsExtractFromPackage(p)
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
	svs, err := StructsExtractFromPackage(path.Join(p, vPrevious))
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

	t, err := templateMigrations(sPrevious, sCurrent, vPrevious, vCurrent)
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

	t, err = templateMigrationsTests(sPrevious, sCurrent, vPrevious, vCurrent)
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

func createDirectoriesForStructs(pathEntities string, structs []StructDef) (paths []string, err error) {
	for _, s := range structs {
		pathEntity := path.Join(pathEntities, s.Name)
		err = DirectoryCreate(pathEntity)
		if err != nil {
			return
		}

		paths = append(paths, pathEntity)
	}
	return
}

func deleteDirectoriesOfNonExistingEntities(pathEntities string, existingEntitiesDirs []string, pathsEntities []string) error {
	for _, d := range existingEntitiesDirs {
		p := path.Join(pathEntities, d)
		s := true
		for _, e := range pathsEntities {
			if p == e {
				s = false
			}
		}
		if s {
			err := DirectoryRemove(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
