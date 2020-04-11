package marketplace

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/marketplace-bl/entities"
)

type ListingsQuery struct {
	Start *ListingID
	Next  int32
}

func (r Resolver) Listings(ctx context.Context, query ListingsQuery) ([]Listing, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return nil, errors.New("Unauthenticated")
	}

	mapListingsToGQL := func(listings []entities.Listing) []Listing {

		gqlListings := []Listing{}

		for _, w := range listings {
			wi := newListing(w, r.listingService, r.communityService)
			gqlListings = append(gqlListings, *wi)
		}

		return gqlListings

	}

	switch query.Start {
	case nil:
		fetchedListings, err := r.listingService.Listings(nil, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapListingsToGQL(fetchedListings), nil
	default:
		id, err := entities.NewListingID(query.Start.String())
		if err != nil {
			return nil, err
		}
		fetchedListings, err := r.listingService.Listings(&id, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapListingsToGQL(fetchedListings), nil

	}

}
