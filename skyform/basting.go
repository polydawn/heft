package skyform

import (
	"fmt"

	sk "github.com/google/skylark"
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
		make(map[string]tlapi.BastingStep),
		make(map[string]tlapi.FormulaContext),
	}}
	switch len(kwargs) {
	case 0: // pass
	default:
		return nil, fmt.Errorf("basting: unexpected keyword arguments")
	}
	switch args.Len() {
	case 0: // pass
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
	default:
		return nil, fmt.Errorf("basting: unexpected extra positional arguments; only 1 is valid")
	}
	return v, err
}

func (x Basting) Type() string          { return "Basting" }
func (x Basting) Truth() sk.Bool        { return true }
func (x Basting) Freeze()               {}                // Freeze is a no-op because we're always a COW structure.
func (x Basting) Hash() (uint32, error) { return 1, nil } // todo
func (x Basting) String() string        { return "<Basting...>" }

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
