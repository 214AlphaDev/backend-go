package wishlist

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/wishlist-bl/entities"
)

type WishesQuery struct {
	Start *UUIDV4
	Next  int32
}

func (r Resolver) Wishes(ctx context.Context, query WishesQuery) ([]Wish, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return nil, errors.New("Unauthenticated")
	}

	mapWishesToGQL := func(wishes []entities.Wish) []Wish {

		gqlWishes := []Wish{}

		for _, w := range wishes {
			wi := newWish(w, r.voteService, r.wishService)
			gqlWishes = append(gqlWishes, *wi)
		}

		return gqlWishes

	}

	switch query.Start {
	case nil:
		fetchedWishes, err := r.wishService.Wishes(nil, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapWishesToGQL(fetchedWishes), nil
	default:
		id, err := entities.NewWishID(query.Start.String())
		if err != nil {
			return nil, err
		}
		fetchedWishes, err := r.wishService.Wishes(&id, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapWishesToGQL(fetchedWishes), nil

	}

}
