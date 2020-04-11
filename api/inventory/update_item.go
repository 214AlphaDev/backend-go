package inventory

import (
	"context"
	"errors"
	"fmt"
	cam "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
	"github.com/214alphadev/inventory-bl/value_objects"
)

type UpdateItemInput struct {
	Item        UUIDV4
	Name        ItemName
	Description ItemDescription
	Story       *ItemStory
	Photo       *Base64ItemPhoto
	Category    string
}

type UpdateItemResponse struct {
	item  *Item
	error error
}

func (r UpdateItemResponse) Error() (*string, error) {
	return toPublicGraphqlError(r.error, "Unauthenticated", "Unauthorized", "ItemDoesNotExist")
}

func (r UpdateItemResponse) Item() *Item {
	return r.item
}

func (r *Resolver) UpdateItem(ctx context.Context, args struct{ Input UpdateItemInput }) (UpdateItemResponse, error) {

	input := args.Input

	memberID := cam.GetAuthenticateMember(ctx)
	if memberID == nil {
		return UpdateItemResponse{error: errors.New("Unauthenticated")}, nil
	}

	member, err := r.communityService.GetMember(*memberID)
	if err != nil {
		return UpdateItemResponse{}, err
	}
	if !member.Admin {
		return UpdateItemResponse{error: errors.New("Unauthorized")}, nil
	}

	itemID, err := entities.NewItemID(input.Item.String())
	if err != nil {
		return UpdateItemResponse{}, err
	}

	item, err := r.itemService.GetByID(itemID)
	if err != nil {
		return UpdateItemResponse{}, err
	}
	if item == nil {
		return UpdateItemResponse{error: errors.New("ItemDoesNotExist")}, nil
	}

	if err := item.ChangeName(input.Name.ItemName); err != nil {
		return UpdateItemResponse{}, err
	}

	if err := item.ChangeDescription(input.Description.ItemDescription); err != nil {
		return UpdateItemResponse{}, err
	}

	switch input.Photo {
	case nil:
		if err := r.itemService.SetPhoto(itemID, nil); err != nil {
			return UpdateItemResponse{}, err
		}
	default:
		if err := r.itemService.SetPhoto(itemID, input.Photo.photo); err != nil {
			return UpdateItemResponse{}, err
		}
	}

	switch input.Story {
	case nil:
		is, err := entities.NewItemStory(nil)
		if err != nil {
			return UpdateItemResponse{}, err
		}
		if err := item.ChangeStory(is); err != nil {
			return UpdateItemResponse{}, err
		}
	default:
		if err := item.ChangeStory(input.Story.ItemStory); err != nil {
			return UpdateItemResponse{}, err
		}
	}

	switch input.Category {
	case "Book":
		item.ChangeCategory(value_objects.Book)
	case "Seed":
		item.ChangeCategory(value_objects.Seed)
	case "Other":
		item.ChangeCategory(value_objects.Other)
	default:
		return UpdateItemResponse{}, fmt.Errorf("received invalid category: %s", input.Category)
	}

	if err := r.itemService.Update(*item); err != nil {
		return UpdateItemResponse{}, err
	}

	return UpdateItemResponse{
		item: newItem(*item, r.itemService),
	}, nil

}
