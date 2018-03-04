package interpret

import (
	"fmt"
	"sort"
	"testing"

	sk "github.com/google/skylark"
)

func TestHello(t *testing.T) {
	script := `iamheft()`
	loader := Loader{}
	globals, err := loader.EvalScript(script)
	fmt.Printf(": %#v\n: %#v\n", globals, err)
	mustGlobalKeys(t, globals, []string{}...)
}

func TestModuleHello(t *testing.T) {
	loader := Loader{
		Psuedofs: map[string]string{
			"fwee.sk": `
def fwee():
	print("kek")
			`,
		},
	}
	script := `load ("fwee.sk", "fwee"); fwee()`
	globals, err := loader.EvalScript(script)
	fmt.Printf(": %#v\n: %#v\n", globals, err)
	mustGlobalKeys(t, globals, "fwee")
}

func TestFormulaFold(t *testing.T) {
	loader := Loader{}
	script := `
f1 = formula({
	"formula":{"action":{
		"exec":["wow", "-c", "as\ndf\n"],
	}},
})
f2 = formula({
	"formula":{"action":{
		"env":{"VAR1":"bees"},
	}},
})
f3 = formula({
	"formula":{"action":{
		"env":{"VAR2":"bats"},
		"exec":["crash", "override"],
	}},
})
f123=f1 + f2 + f3
`
	globals, err := loader.EvalScript(script)
	if err != nil {
		t.Fatalf("error interpreting: %v", err)
	}
	mustGlobalKeys(t, globals, "f1", "f2", "f3", "f123")
	shouldStringEqual(t, `<FormulaUnion:{"formula":{"inputs":{},"action":{"exec":["crash","override"],"env":{"VAR1":"bees","VAR2":"bats"}},"outputs":{}}}>`, globals["f123"].String())
}

func shouldStringEqual(t *testing.T, expect, actual string) {
	t.Helper()
	if expect != actual {
		t.Errorf("want string: %#v\n got string: %#v\n", expect, actual)
	}
}

func mustGlobalKeys(t *testing.T, globals sk.StringDict, wantKeys ...string) {
	t.Helper()
	keys := make([]string, 0, len(globals))
	for k := range globals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sort.Strings(wantKeys)
	keysStr := fmt.Sprintf("%#v", keys)
	wantKeysStr := fmt.Sprintf("%#v", wantKeys)
	if wantKeysStr != keysStr {
		t.Fatalf("want keys: %#v\n got keys: %#v\n", wantKeysStr, keysStr)
	}
}
