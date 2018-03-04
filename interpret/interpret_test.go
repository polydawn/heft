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
