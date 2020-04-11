package inventory

import (
	cd "github.com/214alphadev/community-bl"
	"github.com/214alphadev/inventory-bl/services"
)

type Resolver struct {
	itemService      services.IItemService
	communityService cd.CommunityInterface
}

func newResolver(itemService services.IItemService, communityService cd.CommunityInterface) *Resolver {
	return &Resolver{
		itemService:      itemService,
		communityService: communityService,
	}
}
