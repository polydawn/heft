package commission

import (
	"fmt"
	"strings"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/layout"
)

type (
	// CommissionGraph is a set of catalog/project names mapped to the set of
	// catalogs they import.  The graph also contains nodes which are purely
	// reference to existing released data; these have no imports.
	//
	// Every CatalogName referenced by a CommissionNode's Imports set must
	// be present as another node in the CommissionGraph.  The graph is a
	// DAG; cycles are not permited.
	CommissionGraph map[api.CatalogName]*CommissionNode

	// CommissionNode for each part of the graph recalls which releases are
	// available for this node (e.g. a whole Catalog), and if we have build
	// instructions, it recalls information about that (in which case we're
	// interested only in its imports, at this scale of planning) and the
	// Catalog will contain a dummy release named "candidate" to represent
	// that could-be-built version.
	CommissionNode struct {
		// Catalog of known releases for this node.
		//
		// If we have synthesis instructions, at one of the releases will be
		// named "candidate", and will have a dummy hash value until
		// we actually complete a build for this node.
		//
		// The catalog may also contain several other "real" releases,
		// which already have known hashes and come from previous action.
		Catalog api.Catalog

		// If set, we have build instructions for how to make new stuff at
		// this node, and these are the imports doing so requested.
		// If the ReleaseName in an import is the sentinel value "candidate",
		// that means we depend on the latest build of that other thing;
		// any other imports are of existing built things and thus does not
		// cause any meaningful dependency for execution planning purposes.
		//
		// If CandidateImports is nil, it means there's no basting or build
		// instructions at all for this CatalogName; in this commission,
		// we're purely using already released waypoints for this node.
		CandidateImports map[api.ReleaseItemID]struct{}
	}
)

type CommissionerCfg struct {
	ModuleConfigLoader  layout.Loader
	HitchingInterpreter HitchingInterpreter
}

func (cfg CommissionerCfg) Commission(startAt api.CatalogName, visited CommissionGraph) (CommissionGraph, error) {
	if visited == nil {
		visited = make(CommissionGraph)
	}
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
	visited[startAt] = &CommissionNode{}

	// Load up all the infos for this module.
	moduleCfg, err := cfg.ModuleConfigLoader.LoadModuleConfig(startAt)
	if err != nil {
		return visited, err
	}
	// If this node had other existing releases noted, jot down that info.
	//  It's only interesting if another node references it later, but we can't tell yet.
	if moduleCfg.Catalog != nil {
		visited[startAt].Catalog = *moduleCfg.Catalog
	}
	// If there's no build instructions to consider, then... that's it, return.
	if moduleCfg.HeftScript == "" {
		return visited, nil
	}
	// Interpret the hitching script, then note the imports resulting.
	basting, err := cfg.HitchingInterpreter.Interpret(moduleCfg.HeftScript)
	if err != nil {
		return visited, err
	}
	importSet := projectImportSet(*basting)
	if visited == nil {
		visited = make(CommissionGraph)
	}
	visited[startAt].CandidateImports = importSet

	// We must now recurse through each of these new imports.
	for imp := range importSet {
		visited, err = cfg.commission(imp.CatalogName, visited, backtrace)
		if err != nil {
			return visited, err
		}
	}

	return visited, nil
}

func projectImportSet(basting api.Basting) map[api.ReleaseItemID]struct{} {
	v := make(map[api.ReleaseItemID]struct{})
	for _, step := range basting.Steps {
		for _, imp := range step.Imports {
			v[imp] = struct{}{}
		}
	}
	return v
}

// HitchingInterpreter takes a Hitching script and evaluates it, which is
// expected to yield a single basting.  The interpreter is typically a skylark
// engine, and likely was constructed with some library loading config
// (however, other mocks are used in the commission tests, so that they
// can run without any relationship to the skylark parts of heft).
type HitchingInterpreter interface {
	// REVIEW um do you really want to load the hitching string first?
	// won't that kind of preclude cat?
	// spec out the recursion termination conditions and get back to me.
	Interpret(string) (*api.Basting, error)
}
