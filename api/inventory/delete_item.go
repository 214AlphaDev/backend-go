package inventory

import (
	"context"
	cam "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
)

func (r *Resolver) DeleteItem(ctx context.Context, args struct{ ID UUIDV4 }) (*string, error) {

	memberID := cam.GetAuthenticateMember(ctx)
	if memberID == nil {
		e := "Unauthenticated"
		return &e, nil
	}

	member, err := r.communityService.GetMember(*memberID)
	if err != nil {
		return nil, err
	}
	if !member.Admin {
		e := "Unauthorized"
		return &e, nil
	}

	itemID, err := entities.NewItemID(args.ID.String())
	if err != nil {
		return nil, err
	}

	return toPublicGraphqlError(r.itemService.Delete(itemID), "ItemDoesNotExist")

}
