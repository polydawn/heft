/*
	Skylark bindings.
*/
package skyform

import (
	"bytes"
	"fmt"

	sk "github.com/google/skylark"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	tlapi "go.polydawn.net/go-timeless-api"
)

func MakeFormulaUnion(_ *sk.Thread, _ *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("formula: unexpected positional arguments")
	}
	v := FormulaUnion{}
	// this is, interestingly enough, a job for the obj-obj refmt composition.
	// we'll come back to this after making more patches to hippogryph;
	// this'll be a lot less bothersome to write if we could have map[string]iface.
	// we could also go straight to writing a refmt tokenSource.  skylarkstruct writeJSON has half the work.
	return v, nil
}

var (
	_ sk.Value    = FormulaUnion{}
	_ sk.HasAttrs = FormulaUnion{}
)

type FormulaUnion struct {
	tlapi.FormulaUnion
}

func (s FormulaUnion) Type() string          { return "FormulaUnion" }
func (s FormulaUnion) Truth() sk.Bool        { return true }
func (s FormulaUnion) Freeze()               {}                // Freeze is a no-op because we're always a COW structure.
func (s FormulaUnion) Hash() (uint32, error) { return 1, nil } // todo
func (s FormulaUnion) String() string        { return "<FormulaUnion...>" }

func (s FormulaUnion) Attr(name string) (sk.Value, error) {
	switch name {
	case "toJson":
		return sk.NewBuiltin(name, func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			var buf bytes.Buffer
			if err := refmt.NewMarshallerAtlased(json.EncodeOptions{}, &buf, tlapi.RepeatrAtlas).Marshal(s.FormulaUnion); err != nil {
				return nil, err
			}
			return sk.String(buf.String()), nil
		}), nil
	default:
		return nil, fmt.Errorf("%v has no .%s attribute", s.Type(), name)
	}
}
func (s FormulaUnion) AttrNames() []string {
	return []string{}
}
