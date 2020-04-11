package wishlist

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/wishlist-bl/entities"
	vo "github.com/214alphadev/wishlist-bl/value_objects"
)

type WithdrawVoteResponse struct {
	error     error
	votesLeft *int32
	wish      *Wish
}

func (r WithdrawVoteResponse) Error() (*string, error) {
	switch r.error {
	case nil:
		return nil, nil
	default:
		return toPublicGraphqlError(r.error, "Unauthenticated", "WishDoesNotExist", "NeverVotedOnWish")
	}
}

func (r WithdrawVoteResponse) VotesLeft() *int32 {
	return r.votesLeft
}

func (r WithdrawVoteResponse) Wish() *Wish {
	return r.wish
}

func (r Resolver) WithdrawVote(ctx context.Context, args struct{ Id UUIDV4 }) (WithdrawVoteResponse, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return WithdrawVoteResponse{error: errors.New("Unauthenticated")}, nil
	}

	wishID, err := entities.NewWishID(args.Id.String())
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	memberID, err := vo.NewMemberID(mid.String())
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	err = r.wishService.WithdrawVote(wishID, memberID)
	if err != nil {
		return WithdrawVoteResponse{error: err}, nil
	}

	votesLeft, err := r.voteService.VotesLeft(memberID)
	vl := int32(votesLeft)
	if err != nil {
		return WithdrawVoteResponse{error: err}, nil
	}

	wish, err := r.wishService.GetByID(wishID)
	if err != nil {
		return WithdrawVoteResponse{}, err
	}

	return WithdrawVoteResponse{
		votesLeft: &vl,
		wish:      newWish(*wish, r.voteService, r.wishService),
	}, nil
}
