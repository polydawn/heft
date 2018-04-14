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
}

func (cfg Config) ViewCatalog(
	_ context.Context,
	catalogName api.CatalogName,
) (*api.Catalog, error) {
	moduleInfo, err := cfg.Loader.LoadModuleConfig(catalogName)
	if err != nil {
		return nil, err
	}
	return moduleInfo.Catalog, nil
}
