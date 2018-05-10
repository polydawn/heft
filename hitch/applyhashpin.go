package hitch

import (
	"context"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

func HashPinBasting(
	basting api.Basting,
	hitchViewCatalogTool hitch.ViewCatalog,
) (*api.Basting, error) {
	// wishlist: basting.Clone method, but atm none of our callers care if this gets mutated.
	for _, step := range basting.Steps {
		for key, wareWaypointTuple := range step.Imports {
			catalog, err := hitchViewCatalogTool(context.TODO(), wareWaypointTuple.CatalogName)
			if err != nil {
				return nil, err
			}
			wareID, err := hitch.PluckReleaseItem(*catalog, wareWaypointTuple.ReleaseName, wareWaypointTuple.ItemName)
			if err != nil {
				return nil, err
			}
			step.Formula.Inputs[key] = *wareID
		}
	}
	return &basting, nil
}
