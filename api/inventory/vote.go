package inventory

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
	"github.com/214alphadev/inventory-bl/services"
	vo "github.com/214alphadev/inventory-bl/value_objects"
)

type VoteResponse struct {
	error       error
	item        *entities.ItemID
	itemService services.IItemService
}

func (r VoteResponse) Error() (*string, error) {
	switch r.error {
	case nil:
		return nil, nil
	default:
		return toPublicGraphqlError(r.error, "Unauthenticated", "NoVotesLeft", "ItemDoesNotExist", "AlreadyVoted")
	}
}

func (r VoteResponse) Item() (*Item, error) {

	if r.error != nil {
		return nil, nil
	}

	item, err := r.itemService.GetByID(*r.item)
	if err != nil {
		return nil, err
	}

	return newItem(*item, r.itemService), nil

}

func (r Resolver) Vote(ctx context.Context, args struct{ Id UUIDV4 }) VoteResponse {

	newErrorResponse := func(err error) VoteResponse {
		return VoteResponse{
			error:       err,
			itemService: r.itemService,
		}
	}

	authenticatedMember := GetAuthenticateMember(ctx)
	if authenticatedMember == nil {
		return newErrorResponse(errors.New("Unauthenticated"))
	}

	memberID, err := vo.NewMemberID(authenticatedMember.String())
	if err != nil {
		return newErrorResponse(err)
	}

	itemID, err := entities.NewItemID(args.Id.String())
	if err != nil {
		return newErrorResponse(err)
	}

	if err := r.itemService.Vote(itemID, memberID); err != nil {
		return newErrorResponse(err)
	}

	return VoteResponse{
		itemService: r.itemService,
		item:        &itemID,
	}

}
