/*
	Skylark bindings.
*/
package skyform

import (
	"bytes"
	"fmt"

	sk "github.com/google/skylark"
	"github.com/google/skylark/syntax"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/shared"
	tlapi "go.polydawn.net/go-timeless-api"
)

func MakeFormulaUnion(_ *sk.Thread, _ *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (_ sk.Value, err error) {
	v := FormulaUnion{}
	switch len(kwargs) {
	case 0: // pass
	default:
		return nil, fmt.Errorf("formula: unexpected keyword arguments")
	}
	switch args.Len() {
	case 0: // pass
	case 1: // take this object as a baseline value
		vtoker := NewValueTokenizer()
		vtoker.Bind(args.Index(0))
		umarsh := obj.NewUnmarshaller(tlapi.RepeatrAtlas)
		umarsh.Bind(&v.FormulaUnion)
		pump := shared.TokenPump{
			vtoker,
			umarsh,
		}
		err = pump.Run()
		if err != nil {
			err = fmt.Errorf("formula: could not map object into formula: %s", err)
		}
	default:
		return nil, fmt.Errorf("formula: unexpected extra positional arguments; only 1 is valid")
	}
	return v, err
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
func (s FormulaUnion) String() string        { return "<FormulaUnion:" + s.toJsonString() + ">" }

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
	case "setupHash":
		return sk.NewBuiltin(name, func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			return sk.String(s.Formula.SetupHash()), nil
		}), nil
	default:
		return nil, fmt.Errorf("%v has no .%s attribute", s.Type(), name)
	}
}
func (s FormulaUnion) AttrNames() []string {
	return []string{}
}

func (x FormulaUnion) Binary(op syntax.Token, y sk.Value, side sk.Side) (sk.Value, error) {
	y2, ok := y.(FormulaUnion)
	if !ok {
		return nil, fmt.Errorf("cannot add %s to %s", y.Type(), x.Type())
	}
	if side == sk.Right { // this is terrifying
		x, y2 = y2, x
	}
	switch op {
	case syntax.PLUS:
		z := FormulaUnion{x.FormulaUnion.Apply(y2.FormulaUnion)}
		return z, nil
	default:
		return nil, fmt.Errorf("binary op %q not supported on %s", op, x.Type())
	}
}

func (x FormulaUnion) toJsonString() string {
	var buf bytes.Buffer
	if err := refmt.NewMarshallerAtlased(json.EncodeOptions{Line: []byte{'\n'}, Indent: []byte{'\t'}}, &buf, tlapi.RepeatrAtlas).Marshal(x.FormulaUnion); err != nil {
		panic(err)
	}
	return buf.String()
}
