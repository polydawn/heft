package interpret

import (
	"context"
	"fmt"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
	"go.polydawn.net/heft/skyform"
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
	) (*api.Basting, error)
}

func NewPlannerEvaluator() PlannerEvaluator {
	return plannerEvaluator{
		Interpreter{
			Psuedofs:          heftLibsPsuedofs,
			ModulePredeclared: skyform.AllBuiltins,
		},
	}
}

// This concrete implementation of PlannerEvaluator
//
//	- assumes you're doing heft work, so you want the TL API predeclared;
//	- takes the hitch view funcs and makes them available to the planner;
//	- keeps one interpreter, so library modules will be cached;
//	- does *not* cache the planner script's resulting basting (wrap it);
//
// ... and is otherwise "just" a skylark interpreter.
type plannerEvaluator struct {
	Interpreter
}

const bastingExportName = "pipeline"

func (pe plannerEvaluator) Evaluate(
	ctx context.Context,
	catalogName api.CatalogName,
	plannerScript string,
	hitchViewCatalogTool hitch.ViewCatalog,
) (*api.Basting, error) {
	globals, err := pe.Interpreter.Eval(
		plannerScript,
		string(catalogName),
		nil, // TODO: hitch view injection goes here!
	)
	if err != nil {
		return nil, fmt.Errorf("error generating basting for %q: %s", catalogName, err)
	}
	x, ok := globals[bastingExportName]
	if !ok {
		return nil, fmt.Errorf("error generating basting for %q: expect a variable called %q to be exported!", catalogName, bastingExportName)
	}
	if basting, ok := x.(skyform.Basting); ok {
		return &basting.Basting, nil
	}
	return nil, fmt.Errorf("error generating basting for %q: expect a variable called %q to be a basting object!", catalogName, bastingExportName)
}
