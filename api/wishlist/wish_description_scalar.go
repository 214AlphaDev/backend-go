package wishlist

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/wishlist-bl/entities"
)

type WishDescription struct {
	entities.WishDescription
}

func (WishDescription) ImplementsGraphQLType(name string) bool {
	return name == "WishDescription"
}

func (i *WishDescription) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wd, err := entities.NewWishDescription(*v)
		if err != nil {
			return err
		}
		i.WishDescription = wd
		return nil
	case string:
		wd, err := entities.NewWishDescription(v)
		if err != nil {
			return err
		}
		i.WishDescription = wd
		return nil
	default:
		return errors.New("received invalid type for wish description")
	}

}

func (i WishDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.WishDescription.String())
}
