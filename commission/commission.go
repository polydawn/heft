package commission

import (
	"context"
	"fmt"
	"strings"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/interpret"
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

	// CommissionNode describes an updatable, buildable node in our graph:
	// something we have pipeline generator scripts for, and can evaluate to
	// produce new build plan snapshots with updated releases of its imports.
	//
	// Since commissioning is just about figuring out the graph of dependent
	// builds in a large system, we really only bother to remember the big
	// picture topological information, namely imports.
	// (Despite the fact we had to invoke a planning script to determine those
	// imports, we need not recall the entire basting, because we assume the
	// planner will behave deterministically; thus we may either regenerate or
	// memoize it without loss of generality.
	// Similarly, though many commission nodes may point to one catalogName
	// which is only full of existing releases and has no planner script, thus
	// is not a CommissionNode itself, the loading and caching of that is
	// handled at a lower level than the commissioning logic.)
	CommissionNode struct {
		// If set, we have build instructions for how to make new stuff at
		// this node, and these are the imports doing so requested.
		// If the ReleaseName in an import is the sentinel value "candidate",
		// that means we depend on the latest build of that other thing;
		// any other imports are of existing built things and thus does not
		// cause any meaningful dependency for execution planning purposes.
		//
		// ItemName fields of the import tuple are kept; while they don't
		// determine anything useful about build orders, for leaf nodes that
		// we won't be building, they DO hint which wares we might want to
		// download in advance so we have our whole seed set on hand.
		CandidateImports map[api.ReleaseItemID]struct{}
	}
)

type CommissionerCfg struct {
	ModuleConfigLoader layout.Loader
	HeftInterpreter    interpret.PlannerEvaluator
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
	// If we're here, we expected a commission node.  If there's no planner
	//  script to evaluate, we're in trouble.
	//  (Also, the layout.Loader shouldn't have told us a candidate is possible.)
	if moduleCfg.HeftScript == "" {
		return visited, fmt.Errorf("cannot find planner script for making candidate releases of %q; commission graph broken.", startAt)
	}
	// Interpret the hitching script, then note the imports resulting.
	basting, err := cfg.HeftInterpreter.Evaluate(context.TODO(), startAt, moduleCfg.HeftScript, nil /* TODO */)
	if err != nil {
		return visited, err
	}
	importSet := projectImportSet(*basting)
	if visited == nil {
		visited = make(CommissionGraph)
	}
	visited[startAt].CandidateImports = importSet

	// Review each of the imports.
	//  Any of them where a "candidate" version was selected is something
	//  we now need to recurse through, figuring out what it in turn needs
	//  commissioned so that we can build it.
	for imp := range importSet {
		if imp.ReleaseName != "candidate" {
			continue
		}
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
			// Module-local imports are not interesting at the scale of commission.
			if imp.CatalogName == "wire" {
				continue
			}
			// Keep everything that's not a local wire.
			v[imp] = struct{}{}
		}
	}
	return v
}
