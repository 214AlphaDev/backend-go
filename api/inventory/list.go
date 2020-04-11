package inventory

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
	"github.com/214alphadev/inventory-bl/value_objects"
)

type createItemInput struct {
	Name        ItemName
	Description ItemDescription
	Story       *ItemStory
	Photo       *Base64ItemPhoto
	Category    string
}

type CreateItemResponse struct {
	error error
	item  *Item
}

func (r CreateItemResponse) Error() (*string, error) {
	return toPublicGraphqlError(r.error, "Unauthenticated", "Unauthorized")
}

func (r CreateItemResponse) Item() *Item {
	return r.item
}

func (r *Resolver) Create(ctx context.Context, input createItemInput) (CreateItemResponse, error) {

	mid := community_authentication_middleware.GetAuthenticateMember(ctx)
	if mid == nil {
		return CreateItemResponse{error: errors.New("Unauthenticated")}, nil
	}

	member, err := r.communityService.GetMember(*mid)
	if err != nil {
		return CreateItemResponse{}, err
	}
	if !member.Admin {
		return CreateItemResponse{error: errors.New("Unauthorized")}, nil
	}

	memberID, err := value_objects.NewMemberID(mid.String())
	if err != nil {
		return CreateItemResponse{}, err
	}

	var category value_objects.Category

	switch input.Category {
	case "Book":
		category = value_objects.Book
	case "Seed":
		category = value_objects.Seed
	case "Other":
		category = value_objects.Other
	default:
		return CreateItemResponse{}, fmt.Errorf("received invalid category: %s", input.Category)
	}

	exec := func(itemStory entities.ItemStory) (CreateItemResponse, error) {

		itemID, err := r.itemService.Create(memberID, input.Name.ItemName, input.Description.ItemDescription, itemStory, category)
		if err != nil {
			return CreateItemResponse{}, err
		}

		item, err := r.itemService.GetByID(itemID)
		if err != nil {
			return CreateItemResponse{}, err
		}

		if input.Photo != nil {
			if err := r.itemService.SetPhoto(itemID, input.Photo.Bytes()); err != nil {
				return CreateItemResponse{}, err
			}
		}

		return CreateItemResponse{
			item: newItem(*item, r.itemService),
		}, nil

	}

	switch input.Story {
	case nil:
		story, err := entities.NewItemStory(nil)
		if err != nil {
			return CreateItemResponse{}, err
		}
		return exec(story)
	default:
		return exec(input.Story.ItemStory)
	}

}
