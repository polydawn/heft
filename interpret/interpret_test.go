package interpret

import (
	"sort"
	"testing"

	sk "github.com/google/skylark"
	. "github.com/warpfork/go-wish"
)

func TestHello(t *testing.T) {
	script := `iamheft()`
	loader := Loader{}
	globals, err := loader.EvalScript(script)
	Require(t, err, ShouldEqual, nil)
	Wish(t, globals, ShouldHaveStringDictKeys, []string{})
}

func TestModuleHello(t *testing.T) {
	loader := Loader{
		Psuedofs: map[string]string{
			"fwee.sk": Dedent(`
				def fwee():
					print("kek")
			`),
		},
	}
	script := Dedent(`
		load ("fwee.sk", "fwee")
		fwee()
	`)
	globals, err := loader.EvalScript(script)
	Require(t, err, ShouldEqual, nil)
	Wish(t, globals, ShouldHaveStringDictKeys, []string{"fwee"})
}

func TestFormulaFold(t *testing.T) {
	loader := Loader{}
	script := Dedent(`
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
	`)
	globals, err := loader.EvalScript(script)
	Require(t, err, ShouldEqual, nil)
	Wish(t, globals, ShouldHaveStringDictKeys, []string{"f1", "f2", "f3", "f123"})
	Wish(t, globals["f123"].String(), ShouldEqual, Dedent(`
		<FormulaUnion:{
			"formula": {
				"inputs": {},
				"action": {
					"exec": [
						"crash",
						"override"
					],
					"env": {
						"VAR1": "bees",
						"VAR2": "bats"
					}
				},
				"outputs": {}
			}
		}
		>`))
}

// ShouldHaveStringDictKeys operations on `sk.StringDict` as the 'actual'
// parameter, and a `[]string` as the 'desire' parameter.
func ShouldHaveStringDictKeys(actual, desire interface{}) (string, bool) {
	globals := actual.(sk.StringDict)
	wantKeys := desire.([]string)
	keys := make([]string, 0, len(globals))
	for k := range globals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sort.Strings(wantKeys)
	return ShouldEqual(keys, wantKeys)
}
