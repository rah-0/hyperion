package gen

import (
	"go/format"
	"path"
	"path/filepath"
	"reflect"

	. "github.com/rah-0/hyperion/util"
)

func Generate() error {
	pathEntities, err := filepath.Abs("../entities")
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

func createEntityAtPath(s StructDef, p string) error {
	e := templateEntity(s)
	f, err := format.Source([]byte(e))
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
	err = createEntityAtPath(s, p)
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

	if len(svs) == 0 {
		return false, hv, ErrGeneratorStructNotFound
	} else if len(svs) > 1 {
		return false, hv, ErrGeneratorStructMoreThanOneFound
	}

	sv := svs[0]
	return !reflect.DeepEqual(s, sv), hv, nil
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
				panic("Not implemented") //TODO: create migrator function
			} else {
				err := createEntityAtPath(s, path.Join(p, hv, "entity.go"))
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
