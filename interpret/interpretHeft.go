package interpret

import (
	"context"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

// Interpreting a heft script is composed of preparing a standard skylark interpreter,
// providing it (and all subsequently loaded modules) with the core Timeless API datatypes and manipulators,
// and providing the top level script with the read-only Hitch view functions to a releases DB.
//
// Library modules can all be cached in the usual skylark way;
// the top-level heft scripts only valid result is a basting,
// so those we can cache alone rather than the whole the skylark globals pool.
type PlannerEvaluator interface {
	// Evaluate a planner script, yielding a basting with pinned imports.
	//
	// The catalogName argument is used mainly in logging and error messages, but
	// may also be presumed to be a valid cache key.  Calling Evaluate repeatedly
	// with the same catalogName but varying plannerScript or varying Hitch state
	// is undefined.
	Evaluate(
		ctx context.Context,
		catalogName api.CatalogName,
		plannerScript string,
		hitchViewCatalogTool hitch.ViewCatalog,
	)
}

type plannerEvaluator struct {
	// skylark interpreter, and that's it, isn't it?
}

func (pe plannerEvaluator) Evaluate(
	ctx context.Context,
	catalogName api.CatalogName,
	plannerScript string,
	hitchViewCatalogTool hitch.ViewCatalog,
) (*api.Basting, error) {
	return nil, nil
}
