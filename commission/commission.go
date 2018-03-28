package commission

import (
	"fmt"
	"strings"

	"go.polydawn.net/go-timeless-api"
)

// CommissionGraph is a set of catalog/project names mapped to the set of
// catalogs they import.
//
// REVIEW: it seems entirely possible we need a {catalog,version} tuple
// in the set, because it is still allowed to select non-tip releases,
// and hypothetically even more than one version from the same catalog.
type CommissionGraph map[api.CatalogName]map[api.CatalogName]bool

type CommissionerCfg struct {
	HitchingLoader CommissionTreeViewer
	// todo Interpreter object goes here, interface needed (for mocking too)
}

func (cfg CommissionerCfg) Commission(startAt api.CatalogName, visited CommissionGraph) (CommissionGraph, error) {
	return cfg.commission(startAt, visited, []string{})
}

func (cfg CommissionerCfg) commission(startAt api.CatalogName, visited CommissionGraph, backtrace []string) (CommissionGraph, error) {
	// First, check for cycles.  If this is in our current walk path already, bad.
	nBacktrace := len(backtrace)
	backtrace = append(backtrace, string(startAt))
	for i, backstep := range backtrace[:nBacktrace] {
		if backstep == string(startAt) {
			return visited, fmt.Errorf("cycle found: %s", strings.Join(backtrace[i:nBacktrace], " -> "))
		}
	}

	// If we have visited this before (and not in a cycle), no-op; already have answer.
	if _, ok := visited[startAt]; ok {
		return visited, nil
	}

	// Load up and interpret the hitching script, then note the imports resulting.
	hitching, err := cfg.HitchingLoader.LoadSynthesis(startAt)
	if err != nil {
		return visited, err
	}
	basting, err := cfg.commissionOne(*hitching)
	if err != nil {
		return visited, err
	}
	importSet := projectImportSet(*basting)
	if visited == nil {
		visited = make(CommissionGraph)
	}
	visited[startAt] = importSet

	// We must now recurse through each of these new imports.
	for imp := range importSet {
		visited, err = cfg.commission(imp, visited, backtrace)
		if err != nil {
			return visited, err
		}
	}

	// Now that all the recursing below us is done, pop our backtrace element.
	backtrace = backtrace[:len(backtrace)-1]

	return visited, nil
}

func (cfg CommissionerCfg) commissionOne(Hitching) (*api.Basting, error) {
	// TODO invoke interpreter
	return &api.Basting{}, nil
}

func projectImportSet(basting api.Basting) map[api.CatalogName]bool {
	v := make(map[api.CatalogName]bool)
	for _, step := range basting.Steps {
		for _, imp := range step.Imports {
			v[imp.CatalogName] = true
		}
	}
	return v
}

type Hitching string // a skylark script

type CommissionTreeViewer interface {
	LoadSynthesis(api.CatalogName) (*Hitching, error)
}

type moduleFilesystem interface {
}

var _ moduleFilesystem = mockFs{}

type mockFs map[string]string
