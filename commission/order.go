package commission

import (
	"fmt"
	"sort"

	"go.polydawn.net/go-timeless-api"
)

/*
	Compute a simple topological sort of the steps based on the wire imports.

	We break ties based on lexigraphical sort on the step names.
	We choose this simple tie-breaker rather than attempting any fancier
	logic based on e.g. downstream dependencies, etc, because ease of
	understanding and the simplicity of predicting the result of the sort
	is more important than cleverness; so is the regional stability of the
	sort in the face of changes in other parts of the graph.

	This is almost exactly identical to the toposort used within a
	basting to select its execution order; just at a different scale.
*/
func OrderSteps(graph CommissionGraph) ([]api.CatalogName, error) {
	result := make([]api.CatalogName, 0, len(graph))
	todo := make(map[api.CatalogName]struct{}, len(graph))
	for node := range graph {
		todo[node] = struct{}{}
	}
	edges := []api.CatalogName{}
	for node := range graph {
		edges = append(edges, node)
	}
	sort.Sort(catalogNameByLex(edges))
	for _, node := range edges {
		if err := orderSteps_visit(node, todo, map[api.CatalogName]struct{}{}, &result, graph); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func orderSteps_visit(
	node api.CatalogName,
	todo map[api.CatalogName]struct{},
	loopDetector map[api.CatalogName]struct{},
	result *[]api.CatalogName,
	graph CommissionGraph,
) error {
	// Quick exit if possible.
	if _, ok := todo[node]; !ok {
		return nil
	}
	if _, ok := loopDetector[node]; ok {
		return fmt.Errorf("not a dag: loop detected at %q", node)
	}
	// Mark self for loop detection.
	loopDetector[node] = struct{}{}
	// If this is a non-leaf node, lots of work...
	if graph[node].CandidateImports != nil {
		// Extract any imports which are dependency wiring.
		edges := []api.CatalogName{}
		for edge := range graph[node].CandidateImports {
			// Defense-in-depth sanity check: "wire" are not edge imports, and should not be seen here.
			if edge.CatalogName == "wire" {
				panic("'wire' is not a valid intra-module import in commission graph")
			}
			// Only "candidate" releases mean we need to recurse to build them;
			//  other edges are things we expect as seed set data, and won't build.
			if edge.ReleaseName != "candidate" {
				continue
			}
			// Ok, for candidate releases, we've got a real build dependency edge:
			link := edge.CatalogName
			if _, ok := graph[link]; !ok {
				return fmt.Errorf("invalid commission import: %q imports %q, which has no build plan", node, link)
			}
			edges = append(edges, link)
		}
		// Sort the dependency nodes by name, then recurse.
		//  This sort is necessary for deterministic order of unrelated nodes.
		sort.Sort(catalogNameByLex(edges))
		for _, edge := range edges {
			if err := orderSteps_visit(edge, todo, loopDetector, result, graph); err != nil {
				return nil
			}
		}
	}
	// Done: put this node in the results, and remove from todo set.
	*result = append(*result, node)
	delete(todo, node)
	return nil
}

type catalogNameByLex []api.CatalogName

func (a catalogNameByLex) Len() int           { return len(a) }
func (a catalogNameByLex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a catalogNameByLex) Less(i, j int) bool { return a[i] < a[j] }
