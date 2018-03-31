package layout

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"

	"go.polydawn.net/go-timeless-api"
)

type ModuleConfig struct {
	HeftScript string
	Catalog    *api.Catalog
}

type Loader interface {
	LoadModuleConfig(api.CatalogName) (*ModuleConfig, error)
}

var (
	_ Loader = FSLoader{}
	_ Loader = FixtureLoader{}
)

type FSLoader struct {
	BasePath string
}

const (
	heftPath    = "build.tl"
	catalogPath = "catalog.tl"
)

func (cfg FSLoader) LoadModuleConfig(module api.CatalogName) (*ModuleConfig, error) {
	// First, is there a dir here at all:
	modulePath := filepath.Join(cfg.BasePath, string(module))
	if fi, err := os.Stat(modulePath); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("not a module: %q is not a dir.", module)
	}

	// Start building records.
	mc := ModuleConfig{}
	var any bool

	// Try for a heft script.
	if bs, err := ioutil.ReadFile(filepath.Join(modulePath, heftPath)); err != nil {
		if os.IsNotExist(err) {
			// that's fine; this file is optional.
		} else {
			return nil, fmt.Errorf("unreadable build file: %s", err)
		}
	} else {
		mc.HeftScript = string(bs)
		any = true
	}

	// Try for a catalog snippet.
	if bs, err := ioutil.ReadFile(filepath.Join(modulePath, catalogPath)); err != nil {
		if os.IsNotExist(err) {
			// that's fine; this file is optional.
		} else {
			return nil, fmt.Errorf("unreadable catalog file: %s", err)
		}
	} else {
		if err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bs, &mc.Catalog, api.HitchAtlas); err != nil {
			return nil, fmt.Errorf("failed parsing catalog file for module %q: %s", module, err)
		}
		any = true
	}

	// Final validity check: if we found *none* of the things we recognize:
	if !any {
		return nil, fmt.Errorf("no module config found for %q", module)
	}
	return &mc, nil
}

type FixtureLoader map[api.CatalogName]ModuleConfig

func (cfg FixtureLoader) LoadModuleConfig(module api.CatalogName) (*ModuleConfig, error) {
	if u, ok := cfg[module]; ok {
		return &u, nil
	}
	return nil, fmt.Errorf("404")
}
