package river

import (
	"fmt"
	"os"
	"testing"

	. "github.com/warpfork/go-wish"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/commission"
	"go.polydawn.net/heft/interpret"
	"go.polydawn.net/heft/layout"
	"go.polydawn.net/heft/skyform"
)

type dingusCfg struct {
	CommissionBasePath string
	// The set of project paths asked for on the CLI.
	// This can be selectors, e.g. "foobar.org/..." golang-style.
	// The default behavior is "./...".
	RequestedProjects []string

	// future: sk libs path?
	// future: iam.tl root file that lets you declare a global name prefix for this path?
}

func TestHello(t *testing.T) {
	dingus := dingusCfg{
		"test",
		[]string{"foobar.org/synth"},
	}
	Wish(t, dingus.Bonk(), ShouldEqual, nil)
}

type interpreter struct{}

func (interpreter) Interpret(module api.CatalogName, script string) (*api.Basting, error) {
	voom := &interpret.Interpreter{
		Psuedofs:          LibPsuedofs,
		ModulePredeclared: skyform.AllBuiltins,
	}
	globals, err := voom.Eval(script, string(module), nil)
	if err != nil {
		return nil, fmt.Errorf("error generating basting for %q: %s", module, err)
	}
	x, ok := globals["pipeline"]
	if !ok {
		return nil, fmt.Errorf("error generating basting for %q: expect a variable called %q to be exported!", module, "pipeline")
	}
	if basting, ok := x.(skyform.Basting); ok {
		return &basting.Basting, nil
	}
	return nil, fmt.Errorf("error generating basting for %q: expect a variable called %q to be a basting object!", module, "pipeline")
}

func (cfg dingusCfg) Bonk() error {
	loader := layout.FSLoader{cfg.CommissionBasePath}
	accumulatedGraph := commission.CommissionGraph{}
	commissioner := commission.CommissionerCfg{
		ModuleConfigLoader: loader,
		HeftInterpreter:    interpreter{},
	}

	for _, projPattern := range cfg.RequestedProjects {
		// future: quit ignoring the possibility of "..." patterns.
		module := api.CatalogName(projPattern)
		moduleCfg, err := loader.LoadModuleConfig(module)
		if err != nil {
			// future: this should distinguish for 404's in the case we're on a glob pattern.
			return err
		}
		fmt.Fprintf(os.Stderr, "commissioning for module %q\n", module)
		_ = moduleCfg // it's unfortunate we have to load this to see if we've got something that pattern matches, but then can't hand it to commissioner.

		_, err = commissioner.Commission(module, accumulatedGraph)
		if err != nil {
			return err
		}
	}

	steps, err := commission.OrderSteps(accumulatedGraph)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "total order: %v\n", steps)

	// future: ... evaluate some stuff...?!
	return nil
}
