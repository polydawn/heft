package hitch

import (
	"context"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
	"go.polydawn.net/heft/layout"
)

var (
	_ hitch.ViewCatalog = Config{}.ViewCatalog
)

type Config struct {
	// Loader to use to get catalog info.
	// Typically something that loads from filesystem wrapped in caching.
	//
	// Since it's also loading a commissioning script, this gives us more than
	// we strictly need for our role as a hitch viewer, but that's alright.
	layout.Loader

	builtCandidates map[api.CatalogName]map[api.ItemName]api.WareID
}

func (cfg Config) ViewCatalog(
	_ context.Context,
	catalogName api.CatalogName,
) (*api.Catalog, error) {
	moduleInfo, err := cfg.Loader.LoadModuleConfig(catalogName)
	if err != nil {
		return nil, err
	}
	catalog := moduleInfo.Catalog
	if moduleInfo.HeftScript != "" {
		candidateRelease := api.ReleaseEntry{
			Name:  "candidate",
			Items: nil, // future: some selection funcs will probably want to know what item names are coming!
		}
		catalog.Releases = append([]api.ReleaseEntry{candidateRelease}, catalog.Releases...)
	}
	return moduleInfo.Catalog, nil
}
