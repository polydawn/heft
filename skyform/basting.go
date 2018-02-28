package skyform

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"

	sk "github.com/google/skylark"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj"
	"github.com/polydawn/refmt/shared"
	tlapi "go.polydawn.net/go-timeless-api"
)

type Basting struct {
	tlapi.Basting
}

var (
	_ sk.Value    = Basting{}
	_ sk.HasAttrs = Basting{}
)

func NewBasting(_ *sk.Thread, _ *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (_ sk.Value, err error) {
	v := Basting{tlapi.Basting{
		Steps:    make(map[string]tlapi.BastingStep),
		Exports:  make(map[tlapi.ItemName]tlapi.ReleaseItemID),
		Contexts: make(map[string]tlapi.FormulaContext),
	}}

	// parse kwargs...
	switch len(kwargs) {
	case 0:
		// punt till later: we should have positional arg
	default: // we're going to read every one of em as a step name.
		for _, kw := range kwargs {
			stepName := kw[0].String()
			stepUnion, ok := kw[1].(FormulaUnion) // FIXME: we should probably accept exports here as well.  and make a nice message on cast fail.
			if !ok {
				return nil, fmt.Errorf("basting: expecting kwargs to each be a step definition (which must be a 'formula' type); arg %q was a %s", stepName, kw[1].Type())
			}
			v.Basting.Steps[stepName] = tlapi.BastingStep{
				Imports: stepUnion.Imports,
				Formula: stepUnion.Formula,
			}
			if stepUnion.Context != nil {
				v.Basting.Contexts[stepName] = *stepUnion.Context
			}
		}
		return v, nil
	}
	// ... or, accept a dict as a positional arg and refmt it.
	switch args.Len() {
	case 0:
		return nil, fmt.Errorf("basting: missing positional arguments; one string is required (or, use kwargs)")
	case 1: // take this object as a baseline value
		vtoker := NewValueTokenizer()
		vtoker.Bind(args.Index(0))
		umarsh := obj.NewUnmarshaller(tlapi.HitchAtlas)
		umarsh.Bind(&v.Basting)
		pump := shared.TokenPump{
			vtoker,
			umarsh,
		}
		err = pump.Run()
		return v, nil
	default:
		return nil, fmt.Errorf("basting: unexpected extra positional arguments; only 1 is valid")
	}
}

func (x Basting) Type() string          { return "Basting" }
func (x Basting) Truth() sk.Bool        { return true }
func (x Basting) Freeze()               {}                // Freeze is a no-op because we're always a COW structure.
func (x Basting) Hash() (uint32, error) { return 1, nil } // todo
func (x Basting) String() string        { return "basting(" + x.toJsonString() + ")" }

func (x Basting) Attr(name string) (sk.Value, error) {
	switch name {
	case "addStep":
		return sk.NewBuiltin(name, func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			switch len(kwargs) {
			case 0: // pass
			default:
				return nil, fmt.Errorf("basting: unexpected keyword arguments")
			}
			switch args.Len() {
			case 0: // pass
			default:
				return nil, fmt.Errorf("basting: unexpected positional arguments")
			}
			return nil, nil
		}), nil
	default:
		return nil, fmt.Errorf("%v has no .%s attribute", x.Type(), name)
	}
}

func (x Basting) AttrNames() []string {
	return []string{}
}

func (x Basting) toJsonString() string {
	var buf bytes.Buffer
	if err := refmt.NewMarshallerAtlased(json.EncodeOptions{}, &buf, tlapi.HitchAtlas).Marshal(x.Basting); err != nil {
		panic(err)
	}
	var buf2 bytes.Buffer
	if err := stdjson.Indent(&buf2, buf.Bytes(), "", "\t"); err != nil {
		panic(err)
	}
	return buf2.String()
}
