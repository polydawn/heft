package skyform

import (
	"fmt"
	"os"

	sk "github.com/google/skylark"
)

var AllBuiltins = sk.StringDict{
	"iamheft": sk.NewBuiltin("iamheft", func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
		if thread.Print != nil {
			thread.Print(thread, "yes")
		} else {
			fmt.Fprintln(os.Stderr, "yes")
		}
		return sk.None, nil
	}),
	"formula":       sk.NewBuiltin("formula", MakeFormulaUnion),
	"basting":       sk.NewBuiltin("basting", NewBasting),
	"releaseItemID": sk.NewBuiltin("releaseItemID", NewReleaseItemID),
}
