package inventory

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
	"github.com/214alphadev/inventory-bl/services"
	vo "github.com/214alphadev/inventory-bl/value_objects"
)

type Item struct {
	itemEntity  entities.Item
	itemService services.IItemService
}

func (i Item) ID() (UUIDV4, error) {
	id, err := uuid.FromString(i.itemEntity.ID().String())
	if err != nil {
		return UUIDV4{}, nil
	}
	return UUIDV4{
		UUID: id,
	}, nil
}

func (i Item) Name() ItemName {
	return ItemName{ItemName: i.itemEntity.Name()}
}

func (i Item) Description() ItemDescription {
	return ItemDescription{
		ItemDescription: i.itemEntity.Description(),
	}
}

func (i Item) Story() *ItemStory {
	if i.itemEntity.Story().IsNil() {
		return nil
	}
	return &ItemStory{
		ItemStory: i.itemEntity.Story(),
	}
}

func (i Item) Votes() (int32, error) {
	votes, err := i.itemService.VotesOfItem(i.itemEntity.ID())
	return int32(votes), err
}

func (i Item) VotedOnIt(ctx context.Context) (bool, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return false, errors.New("Unauthenticated")
	}

	memberID, err := vo.NewMemberID(mid.String())
	if err != nil {
		return false, err
	}

	return i.itemService.VotedOnItem(memberID, i.itemEntity.ID())

}

func (i Item) Category() (string, error) {
	switch i.itemEntity.Category() {
	case vo.Book:
		return "Book", nil
	case vo.Seed:
		return "Seed", nil
	case vo.Other:
		return "Other", nil
	default:
		return "", fmt.Errorf("can't handle category: %s", i.itemEntity.Category())
	}
}

func (i Item) Photo() (*Base64ItemPhoto, error) {

	photo, err := i.itemService.GetPhoto(i.itemEntity.ID())
	if err != nil {
		return nil, err
	}

	if photo == nil {
		return nil, nil
	}

	out, err := base64.StdEncoding.DecodeString(photo.String())
	if err != nil {
		return nil, err
	}

	return &Base64ItemPhoto{
		photo: out,
	}, nil

}

func newItem(item entities.Item, itemService services.IItemService) *Item {
	return &Item{
		itemEntity:  item,
		itemService: itemService,
	}
}
