package marketplace

import (
	cd "github.com/214alphadev/community-bl"
	"github.com/214alphadev/marketplace-bl/services"
)

type Resolver struct {
	listingService   services.IListingService
	communityService cd.CommunityInterface
}

func newResolver(listingService services.IListingService, communityService cd.CommunityInterface) *Resolver {
	return &Resolver{
		listingService:   listingService,
		communityService: communityService,
	}
}
