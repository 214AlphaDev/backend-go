package marketplace

import (
	"context"
	"errors"
	cam "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/marketplace-bl/entities"
)

type UpdateListingInput struct {
	ID            ListingID
	Name          ListingName
	Description   ListingDescription
	Photo         Base64ListingPhoto
	PriceCurrency Currency
	PriceAmount   PriceAmount
}

type UpdateListingResponse struct {
	listing *Listing
	error   error
}

func (r UpdateListingResponse) Error() (*string, error) {
	return toPublicGraphqlError(r.error, "Unauthenticated", "Unauthorized", "ListingDoesNotExist")
}

func (r UpdateListingResponse) Listing() *Listing {
	return r.listing
}

func (r *Resolver) UpdateListing(ctx context.Context, args struct{ Input UpdateListingInput }) (UpdateListingResponse, error) {

	input := args.Input

	memberID := cam.GetAuthenticateMember(ctx)
	if memberID == nil {
		return UpdateListingResponse{error: errors.New("Unauthenticated")}, nil
	}

	listingID, err := entities.NewListingID(input.ID.String())
	if err != nil {
		return UpdateListingResponse{}, err
	}

	listing, err := r.listingService.GetByID(listingID)
	if err != nil {
		return UpdateListingResponse{}, err
	}
	if listing == nil {
		return UpdateListingResponse{error: errors.New("ListingDoesNotExist")}, nil
	}

	if memberID.String() != listing.Seller().String() {
		return UpdateListingResponse{error: errors.New("Unauthorized")}, nil
	}

	if err := listing.ChangeName(input.Name.ListingName); err != nil {
		return UpdateListingResponse{}, err
	}

	if err := listing.ChangeDescription(input.Description.ListingDescription); err != nil {
		return UpdateListingResponse{}, err
	}

	if err := listing.ChangePhoto(input.Photo.photo); err != nil {
		return UpdateListingResponse{}, err
	}

	price, err := entities.NewPrice(input.PriceCurrency.currency, input.PriceAmount.amount)
	if err != nil {
		return UpdateListingResponse{}, err
	}
	if err := listing.ChangePrice(price); err != nil {
		return UpdateListingResponse{}, err
	}

	if err := r.listingService.Update(*listing); err != nil {
		return UpdateListingResponse{}, err
	}

	return UpdateListingResponse{
		listing: newListing(*listing, r.listingService, r.communityService),
	}, nil

}
