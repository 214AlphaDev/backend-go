package inventory

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/inventory-bl/entities"
)

type ItemName struct {
	entities.ItemName
}

func (ItemName) ImplementsGraphQLType(name string) bool {
	return name == "ItemName"
}

func (i *ItemName) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wn, err := entities.NewItemName(*v)
		if err != nil {
			return err
		}
		i.ItemName = wn
		return nil
	case string:
		wn, err := entities.NewItemName(v)
		if err != nil {
			return err
		}
		i.ItemName = wn
		return nil
	default:
		return errors.New("received invalid type for item name")
	}

}

func (i ItemName) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ItemName.String())
}
