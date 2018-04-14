package commission

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/interpret"
	"go.polydawn.net/heft/layout"
)

func TestHello(t *testing.T) {
	spore := CommissionerCfg{
		layout.FixtureLoader{
			"foo.org/bar":      layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foible.net/fwoop:v1:osarch"}}))`},
			"foible.net/fwoop": layout.ModuleConfig{Catalog: &api.Catalog{Name: "foible.net/fwoop"}},
		},
		interpret.NewPlannerEvaluator(),
	}
	graph, err := spore.Commission("foo.org/bar", nil)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, graph, ShouldEqual, CommissionGraph{
		"foo.org/bar": &CommissionNode{
			// todo should now have a catalog with 'candidate' release as well
			CandidateImports: map[api.ReleaseItemID]struct{}{
				api.ReleaseItemID{"foible.net/fwoop", "v1", "osarch"}: struct{}{},
			},
		},
	})
}

func TestCycleRejection(t *testing.T) {
	spore := CommissionerCfg{
		layout.FixtureLoader{
			"foo.org/bar":      layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foible.net/fwoop:candidate:osarch", "/2": "foible.net/edge:v1:osarch"}}))`},
			"foible.net/fwoop": layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foible.net/edge:v1:osarch", "/2": "foible.net/feep:candidate:osarch"}}))`},
			"foible.net/feep":  layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foo.org/bar:candidate:osarch"}}))`},
			"foible.net/edge":  layout.ModuleConfig{Catalog: &api.Catalog{Name: "foible.net/edge"}},
		},
		interpret.NewPlannerEvaluator(),
	}
	_, err := spore.Commission("foo.org/bar", nil)
	Wish(t, err, ShouldEqual, fmt.Errorf("cycle found: foo.org/bar -> foible.net/fwoop -> foible.net/feep"))
}

// *exact* same fixture as TestCycleRejection, *except*:
// in this test, our "cycle" actually depends on a concrete version rather than candidate,
// which means it's... not a cycle.
func TestFrakkedCycleAcceptance(t *testing.T) {
	spore := CommissionerCfg{
		layout.FixtureLoader{
			"foo.org/bar":      layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foible.net/fwoop:candidate:osarch", "/2": "foible.net/edge:v1:osarch"}}))`},
			"foible.net/fwoop": layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foible.net/edge:v1:osarch", "/2": "foible.net/feep:candidate:osarch"}}))`},
			"foible.net/feep":  layout.ModuleConfig{HeftScript: `pipeline = basting(stepA=formula({"imports":{"/":"foo.org/bar:v1:osarch"}}))`},
			"foible.net/edge":  layout.ModuleConfig{Catalog: &api.Catalog{Name: "foible.net/edge"}},
		},
		interpret.NewPlannerEvaluator(),
	}
	_, err := spore.Commission("foo.org/bar", nil)
	Wish(t, err, ShouldEqual, nil)
}
