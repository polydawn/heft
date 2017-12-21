package main

import (
	"fmt"
	"os"

	sk "github.com/google/skylark"

	"go.polydawn.net/heft/skyform"
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
		"formula": sk.NewBuiltin("formula", skyform.MakeFormulaUnion),
	}
}

func execfile(filename string) {
	thread := &sk.Thread{}
	globals := newGlobals()
	if err := sk.ExecFile(thread, filename, nil, globals); err != nil {
		fmt.Fprintf(os.Stderr, "larking: %s\n", err)
		os.Exit(4)
	}
}
