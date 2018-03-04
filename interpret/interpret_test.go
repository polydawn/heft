package interpret

import (
	"fmt"
	"testing"
)

func TestHello(t *testing.T) {
	script := `iamheft()`
	loader := Loader{}
	globals, err := loader.EvalScript(script)
	fmt.Printf(": %#v\n: %#v\n", globals, err)
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
}
