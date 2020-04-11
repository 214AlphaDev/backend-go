package inventory

import (
	"context"
	"errors"
	. "github.com/214alphadev/community-authentication-middleware"
	"github.com/214alphadev/inventory-bl/entities"
)

type ItemsQuery struct {
	Start *UUIDV4
	Next  int32
}

func (r Resolver) Items(ctx context.Context, query ItemsQuery) ([]Item, error) {

	mid := GetAuthenticateMember(ctx)
	if mid == nil {
		return nil, errors.New("Unauthenticated")
	}

	mapItemsToGQL := func(items []entities.Item) []Item {

		gqlItems := []Item{}

		for _, w := range items {
			wi := newItem(w, r.itemService)
			gqlItems = append(gqlItems, *wi)
		}

		return gqlItems

	}

	switch query.Start {
	case nil:
		fetchedItems, err := r.itemService.Items(nil, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapItemsToGQL(fetchedItems), nil
	default:
		id, err := entities.NewItemID(query.Start.String())
		if err != nil {
			return nil, err
		}
		fetchedItems, err := r.itemService.Items(&id, uint32(query.Next))
		if err != nil {
			return nil, err
		}
		return mapItemsToGQL(fetchedItems), nil

	}

}
