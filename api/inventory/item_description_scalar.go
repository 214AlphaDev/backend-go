package inventory

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/inventory-bl/entities"
)

type ItemDescription struct {
	entities.ItemDescription
}

func (ItemDescription) ImplementsGraphQLType(name string) bool {
	return name == "ItemDescription"
}

func (i *ItemDescription) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wd, err := entities.NewItemDescription(*v)
		if err != nil {
			return err
		}
		i.ItemDescription = wd
		return nil
	case string:
		wd, err := entities.NewItemDescription(v)
		if err != nil {
			return err
		}
		i.ItemDescription = wd
		return nil
	default:
		return errors.New("received invalid type for item description")
	}

}

func (i ItemDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ItemDescription.String())
}
