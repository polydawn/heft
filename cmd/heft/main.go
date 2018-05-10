package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/commission"
	"go.polydawn.net/heft/interpret"
	"go.polydawn.net/heft/layout"
)

func main() {
	if err := Main(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(4)
	}
}

func Main(requestedProjects ...string) error {
	loader := layout.FSLoader{os.Getenv("HEFT_ROOT")}
	accumulatedGraph := commission.CommissionGraph{}
	commissioner := commission.CommissionerCfg{
		ModuleConfigLoader: loader,
		HeftInterpreter:    interpret.NewPlannerEvaluator(),
	}

	for _, projPattern := range requestedProjects {
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

	projectsOrdering, err := commission.OrderSteps(accumulatedGraph)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "total order: %v\n\n", projectsOrdering)

	for _, projectName := range projectsOrdering {
		// eval planner (again; should be same)
		fmt.Fprintf(os.Stderr, "evaluing project %q\n\n", projectName)
		moduleCfg, err := loader.LoadModuleConfig(projectName)
		if err != nil {
			return err
		}
		basting, err := commissioner.HeftInterpreter.Evaluate(
			context.Background(),
			projectName,
			moduleCfg.HeftScript,
			nil,
		)
		if err != nil {
			return err
		}

		// Resolve basting's imports to full on hashes.
		// TODO -- should just need a hitch plus its dict of builtCandidates wip results.

		// exec!
		cmd := exec.Command("repeatr", "batch", "/dev/fd/0")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		in, _ := cmd.StdinPipe()
		if err := refmt.NewMarshallerAtlased(json.EncodeOptions{}, in, api.HitchAtlas).Marshal(basting); err != nil {
			return err
		}
		if err := cmd.Run(); err != nil {
			return err
		}

		// Parse that resultgroup.
		//  These are now items for the "candidate" release -- add them to our builtCandidates set.
		// TODO
	}
	return nil
}
