package commission

import (
	"fmt"
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

func (mockInterpreter) Interpret(Hitching) (*api.Basting, error) {
	// this still needs to select a *releasename*, so
	//  it will still need to be constructed with a "hitch" view caller.
	//  actually that should be part of the params on the interface
	//   because it's both standard and may chnage behavior from step to step.

	// yes we really should store the releasename in the visit map.
	//  references to not-candidate aren't going to cause execution; but they
	//   also shouldn't render the same way in the graphviz.
	//   (unclear if it should count or not for cyclicy.

	return &api.Basting{}, nil
}

func TestHello(t *testing.T) {
	spore := CommissionerCfg{
		mockLoader{
			"foo.org/bar": hitchingNode("um basting plz"),
		},
		mockInterpreter{},
	}
	graph, err := spore.Commission("foo.org/bar", nil)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, graph, ShouldEqual, CommissionGraph{
		"foo.org/bar": map[api.CatalogName]bool{},
	})
}
