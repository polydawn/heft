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
