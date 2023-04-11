package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
)

const mapImagesFolder string = "./../../../frontend/public/images/maps"

func TestMapImagesComplete(t *testing.T) {
	mapNames, err := getMapNames()
	if err != nil {
		t.Fatalf("could not determine map names: %v", err)
	}

	for _, mapName := range mapNames {
		t.Run(mapName, func(t *testing.T) {
			filename := mapName + ".jpg"
			expectedPath, err := filepath.Abs(filepath.Join(mapImagesFolder, filename))
			if err != nil {
				t.Errorf("could not determine image path: %v", err)
				return
			}
			var info fs.FileInfo
			if info, err = os.Stat(expectedPath); (err != nil && errors.Is(err, os.ErrNotExist)) || info.Name() != filename {
				t.Errorf("expected image file %s for %s not found", expectedPath, mapName)
			}
		})
	}
}

const mapTypeName string = "Map"

func getMapNames() (mapNames []string, err error) {
	cfg := &packages.Config{Mode: packages.NeedTypes}
	var pkgs []*packages.Package
	pkgs, err = packages.Load(cfg, "pattern=github.com/redraskal/r6-dissect/dissect")
	if err != nil {
		return
	}

	if packages.PrintErrors(pkgs) > 0 {
		err = errors.New("got error loading package")
		return
	}

	if len(pkgs) == 0 {
		err = errors.New("found no packages")
		return
	}
	scope := pkgs[0].Types.Scope()
	mapType := scope.Lookup(mapTypeName).Type()

	for _, varName := range scope.Names() {
		obj := scope.Lookup(varName)
		if obj.Type() == mapType && varName != mapTypeName {
			mapNames = append(mapNames, varName)
		}
	}
	return
}
