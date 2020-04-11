package inventory

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/inventory-bl/entities"
)

type ItemStory struct {
	entities.ItemStory
}

func (ItemStory) ImplementsGraphQLType(name string) bool {
	return name == "ItemStory"
}

func (i *ItemStory) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		ws, err := entities.NewItemStory(v)
		if err != nil {
			return err
		}
		i.ItemStory = ws
		return nil
	case string:
		ws, err := entities.NewItemStory(&v)
		if err != nil {
			return err
		}
		i.ItemStory = ws
		return nil
	default:
		return errors.New("received invalid type for item story")
	}

}

func (i ItemStory) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ItemStory.String())
}
