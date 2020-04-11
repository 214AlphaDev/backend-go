package marketplace

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/satori/go.uuid"
	cd "github.com/214alphadev/community-bl"
	"github.com/214alphadev/marketplace-bl/entities"
	"github.com/214alphadev/marketplace-bl/services"
)

type Listing struct {
	listingEntity    entities.Listing
	listingService   services.IListingService
	communityService cd.CommunityInterface
}

func (i Listing) ID() (ListingID, error) {
	return ListingID{id: i.listingEntity.ID()}, nil
}

func (i Listing) Name() ListingName {
	return ListingName{ListingName: i.listingEntity.Name()}
}

func (i Listing) Description() ListingDescription {
	return ListingDescription{
		ListingDescription: i.listingEntity.Description(),
	}
}

func (i Listing) Price() Price {
	return newPrice(i.listingEntity.Price())
}

func (i Listing) Photo() Base64ListingPhoto {
	return Base64ListingPhoto{photo: i.listingEntity.Photo()}
}

func (i Listing) ListedAt() graphql.Time {
	return graphql.Time{
		Time: i.listingEntity.CreatedAt().Time(),
	}
}

func (i Listing) Seller() (Seller, error) {

	memberId, err := uuid.FromString(i.listingEntity.Seller().String())
	if err != nil {
		return Seller{}, err
	}

	member, err := i.communityService.GetMember(memberId)
	if err != nil {
		return Seller{}, err
	}

	return newSeller(member), nil
}

func newListing(listing entities.Listing, listingService services.IListingService, communityService cd.CommunityInterface) *Listing {
	return &Listing{
		listingEntity:    listing,
		listingService:   listingService,
		communityService: communityService,
	}
}
