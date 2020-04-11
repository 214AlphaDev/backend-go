package inventory

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
	vo "github.com/214alphadev/inventory-bl/value_objects"
)

type WithdrawVoteResponse struct {
	error error
	item  *Item
}

func (r WithdrawVoteResponse) Error() (*string, error) {
	switch r.error {
	case nil:
		return nil, nil
	default:
		return toPublicGraphqlError(r.error, "Unauthenticated", "ItemDoesNotExist", "NeverVotedOnItem")
	}
}

func (r WithdrawVoteResponse) Item() *Item {
	return r.item
}

func (r Resolver) WithdrawVote(ctx context.Context, args struct{ Id UUIDV4 }) (WithdrawVoteResponse, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return WithdrawVoteResponse{error: errors.New("Unauthenticated")}, nil
	}

	itemID, err := entities.NewItemID(args.Id.String())
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	memberID, err := vo.NewMemberID(mid.String())
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	err = r.itemService.WithdrawVote(itemID, memberID)
	if err != nil {
		return WithdrawVoteResponse{error: err}, nil
	}

	item, err := r.itemService.GetByID(itemID)
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	return WithdrawVoteResponse{
		item: newItem(*item, r.itemService),
	}, nil
}
