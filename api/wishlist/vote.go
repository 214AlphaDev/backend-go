package wishlist

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/wishlist-bl/entities"
	. "github.com/214alphadev/wishlist-bl/services"
	. "github.com/214alphadev/wishlist-bl/value_objects"
)

type VoteResponse struct {
	error       error
	wish        *entities.WishID
	wishService IWishService
	voteService IVoteService
}

func (r VoteResponse) Error() (*string, error) {
	switch r.error {
	case nil:
		return nil, nil
	default:
		return toPublicGraphqlError(r.error, "Unauthenticated", "NoVotesLeft", "WishDoesNotExist", "AlreadyVoted")
	}
}

func (r VoteResponse) VotesLeft(ctx context.Context) (*int32, error) {

	authenticatedMember := GetAuthenticateMember(ctx)
	if authenticatedMember == nil {
		return nil, errors.New("member is not authenticated")
	}

	memberID, err := NewMemberID(authenticatedMember.String())
	if err != nil {
		return nil, err
	}

	leftVotes, err := r.voteService.VotesLeft(memberID)
	if err != nil {
		return nil, err
	}

	leftVotesI := int32(leftVotes)

	return &leftVotesI, nil

}

func (r VoteResponse) Wish() (*Wish, error) {

	if r.error != nil {
		return nil, nil
	}

	wish, err := r.wishService.GetByID(*r.wish)
	if err != nil {
		return nil, err
	}

	return newWish(*wish, r.voteService, r.wishService), nil

}

func (r Resolver) Vote(ctx context.Context, args struct{ Id UUIDV4 }) VoteResponse {

	newErrorResponse := func(err error) VoteResponse {
		return VoteResponse{
			error:       err,
			wishService: r.wishService,
			voteService: r.voteService,
		}
	}

	authenticatedMember := GetAuthenticateMember(ctx)
	if authenticatedMember == nil {
		return newErrorResponse(errors.New("Unauthenticated"))
	}

	memberID, err := NewMemberID(authenticatedMember.String())
	if err != nil {
		return newErrorResponse(err)
	}

	wishID, err := entities.NewWishID(args.Id.String())
	if err != nil {
		return newErrorResponse(err)
	}

	if err := r.wishService.Vote(wishID, memberID); err != nil {
		return newErrorResponse(err)
	}

	return VoteResponse{
		wishService: r.wishService,
		voteService: r.voteService,
		wish:        &wishID,
	}

}
