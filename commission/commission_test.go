package commission

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	"go.polydawn.net/go-timeless-api"
)

type mockLoader map[string]CommissionStepUnion

func (m mockLoader) LoadSynthesis(catname api.CatalogName) (*CommissionStepUnion, error) {
	if u, ok := m[string(catname)]; ok {
		return &u, nil
	}
	return nil, fmt.Errorf("404")
}

func hitchingNode(content string) CommissionStepUnion {
	cast := Hitching(content)
	return CommissionStepUnion{Hitching: &cast}
}

type mockInterpreter struct{}

func (mockInterpreter) Interpret(script Hitching) (*api.Basting, error) {
	// this still needs to select a *releasename*, so
	//  it will still need to be constructed with a "hitch" view caller.
	//  actually that should be part of the params on the interface
	//   because it's both standard and may chnage behavior from step to step.

	// yes we really should store the releasename in the visit map.
	//  references to not-candidate aren't going to cause execution; but they
	//   also shouldn't render the same way in the graphviz.
	//   (unclear if it should count or not for cyclicy.

	split := strings.Split(string(script), ",")
	imps := make(map[api.AbsPath]api.ReleaseItemID, len(split))
	for i, hunk := range split {
		imp, err := api.ParseReleaseItemID(hunk)
		if err != nil {
			panic(err) // your fixture is wrong
		}
		imp.ReleaseName = "v1" // todo see above comments about needing catalog viewer
		imps[api.AbsPath(fmt.Sprintf("/%d", i))] = imp
	}
	return &api.Basting{
		Steps: map[string]api.BastingStep{
			"astep": {
				Imports: imps,
				// a real basting would have formula, etc here, but for
				//  what we're testing in this package none of that is relevant.
			},
		},
	}, nil
}

func TestHello(t *testing.T) {
	spore := CommissionerCfg{
		mockLoader{
			"foo.org/bar":      hitchingNode("foible.net/fwoop"),
			"foible.net/fwoop": CommissionStepUnion{Catalog: &api.Catalog{Name: "foible.net/fwoop"}},
		},
		mockInterpreter{},
	}
	graph, err := spore.Commission("foo.org/bar", nil)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, graph, ShouldEqual, CommissionGraph{
		"foo.org/bar": map[api.CatalogName]bool{
			"foible.net/fwoop": true,
		},
		"foible.net/fwoop": nil,
	})
}
