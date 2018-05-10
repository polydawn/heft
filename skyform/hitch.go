package skyform

import (
	"fmt"

	sk "github.com/google/skylark"

	"go.polydawn.net/go-timeless-api/hitch"
)

type Hitch struct {
	viewCatalogTool hitch.ViewCatalog
}

var (
	_ sk.Value    = Hitch{}
	_ sk.HasAttrs = Hitch{}
)

func HitchSingleton(hitchViewCatalogTool hitch.ViewCatalog) sk.Value {
	v := sk.NewBuiltin("hitch", func(_ *sk.Thread, _ *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (_ sk.Value, err error) {
		return Hitch{hitchViewCatalogTool}, nil
	})
	return v
}

func (x Hitch) Type() string          { return "Hitch" }
func (x Hitch) Truth() sk.Bool        { return true }
func (x Hitch) Freeze()               {}
func (x Hitch) Hash() (uint32, error) { return 1, nil }
func (x Hitch) String() string        { return "hitch([unrepresentable])" }

func (x Hitch) Attr(name string) (sk.Value, error) {
	switch name {
	case "viewCatalog":
		// TODO ... enough new types to translate to skyform that i'm seriously considering codegen
		return nil, nil
	default:
		return nil, fmt.Errorf("%v has no .%s attribute", x.Type(), name)
	}
}

func (x Hitch) AttrNames() []string {
	return []string{}
}
