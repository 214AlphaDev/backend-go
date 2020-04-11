package wishlist

import (
	"context"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/wishlist-bl/entities"
	"github.com/214alphadev/wishlist-bl/services"
	vo "github.com/214alphadev/wishlist-bl/value_objects"
)

type Wish struct {
	wishEntity  entities.Wish
	voteService services.IVoteService
	wishService services.IWishService
}

func (w Wish) ID() (UUIDV4, error) {
	id, err := uuid.FromString(w.wishEntity.ID().String())
	if err != nil {
		return UUIDV4{}, nil
	}
	return UUIDV4{
		UUID: id,
	}, nil
}

func (w Wish) Name() WishName {
	return WishName{WishName: w.wishEntity.Name()}
}

func (w Wish) Description() WishDescription {
	return WishDescription{
		WishDescription: w.wishEntity.Description(),
	}
}

func (w Wish) Story() *WishStory {
	if w.wishEntity.Story().IsNil() {
		return nil
	}
	return &WishStory{
		WishStory: w.wishEntity.Story(),
	}
}

func (w Wish) Votes() (int32, error) {
	votes, err := w.wishService.VotesOfWish(w.wishEntity.ID())
	return int32(votes), err
}

func (w Wish) Category() (string, error) {
	switch w.wishEntity.Category() {
	case vo.Seed:
		return "Seed", nil
	case vo.Book:
		return "Book", nil
	case vo.Other:
		return "Other", nil
	default:
		return "", fmt.Errorf("received invalid category: %s", w.wishEntity.Category())
	}
}

func (w Wish) VotedOnIt(ctx context.Context) (bool, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return false, errors.New("Unauthenticated")
	}

	memberID, err := vo.NewMemberID(mid.String())
	if err != nil {
		return false, err
	}

	votes, err := w.voteService.Votes(memberID)
	if err != nil {
		return false, err
	}

	for _, v := range votes {
		if v.WishID().String() == w.wishEntity.ID().String() {
			return true, nil
		}
	}

	return false, nil
}

func newWish(wish entities.Wish, voteService services.IVoteService, wishService services.IWishService) *Wish {
	return &Wish{
		wishEntity:  wish,
		voteService: voteService,
		wishService: wishService,
	}
}
