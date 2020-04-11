package wishlist

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/wishlist-bl/entities"
	"github.com/214alphadev/wishlist-bl/value_objects"
	. "github.com/214alphadev/wishlist-bl/value_objects"
)

type createWishInput struct {
	Name        WishName
	Description WishDescription
	Story       *WishStory
	Category    string
}

type CreateWishResponse struct {
	error error
	wish  *Wish
}

func (r CreateWishResponse) Error() *string {
	switch r.error {
	case nil:
		return nil
	default:
		e := r.error.Error()
		return &e
	}
}

func (r CreateWishResponse) Wish() *Wish {
	return r.wish
}

func (r *Resolver) Create(ctx context.Context, input createWishInput) (CreateWishResponse, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return CreateWishResponse{error: errors.New("Unauthenticated")}, nil
	}

	memberID, err := NewMemberID(mid.String())
	if err != nil {
		return CreateWishResponse{}, err
	}

	var category value_objects.Category
	switch input.Category {
	case "Seed":
		category = Seed
	case "Book":
		category = Book
	case "Other":
		category = Other
	default:
		return CreateWishResponse{}, fmt.Errorf(`can't handle category: "%s"'`, input.Category)
	}

	exec := func(wishStory entities.WishStory) (CreateWishResponse, error) {

		wishID, err := r.wishService.Create(memberID, input.Name.WishName, input.Description.WishDescription, wishStory, category)
		if err != nil {
			return CreateWishResponse{}, err
		}

		wish, err := r.wishService.GetByID(wishID)
		if err != nil {
			return CreateWishResponse{}, err
		}

		return CreateWishResponse{
			wish: newWish(*wish, r.voteService, r.wishService),
		}, nil

	}

	switch input.Story {
	case nil:
		story, err := entities.NewWishStory(nil)
		if err != nil {
			return CreateWishResponse{}, err
		}
		return exec(story)
	default:
		return exec(input.Story.WishStory)
	}

}
