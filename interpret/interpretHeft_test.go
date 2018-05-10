package interpret

// These tests are a little different than the baseline interpreter tests
// because they A) return bastings, and B) need a hitch DB handle.
//
// So, we can use these to exercise release selection func defs in skylark.

import (
	"context"
	"testing"

	. "github.com/warpfork/go-wish"

	"go.polydawn.net/go-timeless-api"
)

func TestHeftHitching(t *testing.T) {
	script := Dedent(`
		pipeline = basting({})
	`)

	pe := NewPlannerEvaluator()
	basting, err := pe.Evaluate(
		context.Background(),
		api.CatalogName("foo.bar/baz"),
		script,
		nil, // TODO: make hitch lens from layout.FixtureLoader
	)

	Wish(t, err, ShouldEqual, nil)
	Wish(t, basting, ShouldEqual, &api.Basting{
		Steps:    map[string]api.BastingStep{},
		Exports:  map[api.ItemName]api.ReleaseItemID{},
		Contexts: map[string]api.FormulaContext{},
	})
}
