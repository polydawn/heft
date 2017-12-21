package main

import (
	"fmt"
	"os"

	sk "github.com/google/skylark"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	tlapi "go.polydawn.net/go-timeless-api"
)

func main() {
	execfile(os.Args[1])
}

// Makes a new globals dict with our favorite custom bits in it.
func newGlobals() sk.StringDict {
	return sk.StringDict{
		"iamheft": sk.NewBuiltin("iamheft", func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			if thread.Print != nil {
				thread.Print(thread, "yes")
			} else {
				fmt.Fprintln(os.Stderr, "yes")
			}
			return sk.None, nil
		}),
		"writeFormula": sk.NewBuiltin("writeFormula", func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			var frm tlapi.FormulaUnion
			if err := refmt.NewMarshallerAtlased(json.EncodeOptions{}, os.Stdout, tlapi.RepeatrAtlas).Marshal(frm); err != nil {
				return nil, err
			}

			return sk.None, nil
		}),
	}
}

func execfile(filename string) {
	thread := &sk.Thread{}
	globals := newGlobals()
	if err := sk.ExecFile(thread, filename, nil, globals); err != nil {
		fmt.Fprintf(os.Stderr, "larking: %s", err)
		os.Exit(4)
	}
}
