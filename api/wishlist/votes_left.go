package wishlist

import (
	"context"
	. "github.com/214alphadev/community-authentication-middleware"
	. "github.com/214alphadev/wishlist-bl/value_objects"
)

func (r Resolver) VotesLeft(ctx context.Context) (*int32, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return nil, nil
	}

	memberID, err := NewMemberID(mid.String())
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
