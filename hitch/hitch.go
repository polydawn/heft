package hitch

import (
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/heft/commission"
)

type View interface {
	Catalog(api.CatalogName) (*api.Catalog, error)
}

// need to:
// - look up a catalog file (and yes we might not've recursed to that node yet in commission mode
// - oh, right.  recurse to that node in commission mode, because we need to know if to make a candidate.

type CommissionView struct {
	Commission commission.CommissionGraph
}

func (cv CommissionView) Catalog(api.CatalogName) (*api.Catalog, error) {
	// we *do* want to load here.
	// it should already be loaded; we don't have to push cache back up and around.
}

// so yeah also this interface is here to stay even if it wraps outside exec api later.

// ... so who's job is that, exactly?
// - If we have a freestanding tool for release/waypoint grooming, that doesn't want to know about this.  It's just a commissioning thing.
// - If we're assuming that any process can be a planner, and that api is a CLI... then we *have* to... hint to it with env vars?!
//   - Are there any other APIs we can use?  Planner has to ask, so, inb4 just hand down data; that doesn't work.api

// You're also a tad fucked in terms of *how much* data you have to hand down.
// Given that we don't know in advance what the planner is going to ask for (that's the point, after all),
// we can't hand down the entire mass of obtw-these-have-candidates-with-this-metadata in env vars, that's for sure;
// doubly impossible to do so *for the entire universe of catalogs we might be allowed to reference* at once, which is what would be required.
// So... do we have any choice at all except to each hitch to read commission node candidate config? >:/
// We could put the search pattern in an env var, I guess.  That seems harmless enough.
