package skyform

import (
	"fmt"

	sk "github.com/google/skylark"
	tlapi "go.polydawn.net/go-timeless-api"
)

type ReleaseItemID struct {
	tlapi.ReleaseItemID
}

var (
	_ sk.Value    = ReleaseItemID{}
	_ sk.HasAttrs = ReleaseItemID{}
)

/*
	Constructor calls should use either kwargs or a single string which will be parsed.
	The version can be omitted in kwargs mode (because that partial specification is
	useful as an argument to version lookups that peek at hitch db state).

	Most functions that expect a ReleaseItemID should also try to accept a string
	and do the parse implicitly.
*/
func NewReleaseItemID(_ *sk.Thread, _ *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (_ sk.Value, err error) {
	switch len(kwargs) {
	case 0:
		// punt till later: we should have positional arg
	case 2, 3:
		v := ReleaseItemID{}
		for _, kw := range kwargs {
			key, _ := sk.AsString(kw[0])
			val, ok := sk.AsString(kw[1])
			if !ok {
				return nil, fmt.Errorf("releaseItemID: unexpected keyword arguments -- values must all be strings")
			}
			switch key {
			case "catalog":
				v.CatalogName = tlapi.CatalogName(val)
			case "version":
				v.ReleaseName = tlapi.ReleaseName(val)
			case "item":
				v.ItemName = tlapi.ItemName(val)
			default:
				return nil, fmt.Errorf("releaseItemID: unexpected keyword arguments -- %q is not one of ['catalog', 'version', 'item']", kw[0].String())
			}
		}
		if v.CatalogName == "" || v.ItemName == "" {
			fmt.Errorf("releaseItemID: missing required keyword arguments -- at least 'catalog' and 'item' are required")
		}
		return v, nil
	default:
		return nil, fmt.Errorf("releaseItemID: unexpected keyword arguments -- 'catalog', 'version', and 'item' are the recognized keys")
	}
	switch args.Len() {
	case 0:
		return nil, fmt.Errorf("releaseItemID: missing positional arguments; one string is required (or, use kwargs)")
	case 1: // parse it
		v, err := tlapi.ParseReleaseItemID(args[0].String())
		return ReleaseItemID{v}, err
	default:
		return nil, fmt.Errorf("releaseItemID: unexpected extra positional arguments; only 1 is valid")
	}
}

func (x ReleaseItemID) asValue() sk.Value { return sk.String(x.String()) }

func (x ReleaseItemID) Type() string          { return "ReleaseItemID" }
func (x ReleaseItemID) Truth() sk.Bool        { return true }
func (x ReleaseItemID) Freeze()               {}                // Freeze is a no-op because we're always a COW structure.
func (x ReleaseItemID) Hash() (uint32, error) { return 1, nil } // todo
func (x ReleaseItemID) String() string        { return x.ReleaseItemID.String() }

func (x ReleaseItemID) Attr(name string) (sk.Value, error) {
	switch name {
	case "catalog":
		return sk.String(x.ReleaseItemID.CatalogName), nil
	case "version":
		return sk.String(x.ReleaseItemID.ReleaseName), nil
	case "item":
		return sk.String(x.ReleaseItemID.ItemName), nil
	default:
		return nil, fmt.Errorf("%v has no .%s attribute", x.Type(), name)
	}
}

func (x ReleaseItemID) AttrNames() []string {
	return []string{
		"catalog",
		"version",
		"item",
	}
}
