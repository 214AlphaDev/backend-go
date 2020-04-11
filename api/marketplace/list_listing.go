package marketplace

import (
	"context"
	"github.com/pkg/errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/marketplace-bl/entities"
	. "github.com/214alphadev/marketplace-bl/value_objects"
)

type createListingInput struct {
	Name          ListingName
	Description   ListingDescription
	Photo         Base64ListingPhoto
	PriceCurrency Currency
	PriceAmount   PriceAmount
}

type CreateListingResponse struct {
	error   error
	listing *Listing
}

func (r CreateListingResponse) Error() (*string, error) {
	return toPublicGraphqlError(r.error, "Unauthenticated", "Unauthorized")
}

func (r CreateListingResponse) Listing() *Listing {
	return r.listing
}

func (r *Resolver) ListNewListing(ctx context.Context, input createListingInput) (CreateListingResponse, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return CreateListingResponse{error: errors.New("Unauthenticated")}, nil
	}

	memberID, err := NewMemberID(mid.String())
	if err != nil {
		return CreateListingResponse{}, err
	}

	price, err := entities.NewPrice(input.PriceCurrency.String(), input.PriceAmount.amount)
	if err != nil {
		return CreateListingResponse{}, err
	}

	listingID, err := r.listingService.List(memberID, input.Name.ListingName, input.Description.ListingDescription, price, input.Photo.photo)
	if err != nil {
		return CreateListingResponse{}, err
	}

	listing, err := r.listingService.GetByID(listingID)
	if err != nil {
		return CreateListingResponse{}, err
	}

	return CreateListingResponse{listing: newListing(*listing, r.listingService, r.communityService)}, nil

}
