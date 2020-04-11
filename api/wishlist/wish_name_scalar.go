package wishlist

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/wishlist-bl/entities"
)

type WishName struct {
	entities.WishName
}

func (WishName) ImplementsGraphQLType(name string) bool {
	return name == "WishName"
}

func (i *WishName) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wn, err := entities.NewWishName(*v)
		if err != nil {
			return err
		}
		i.WishName = wn
		return nil
	case string:
		wn, err := entities.NewWishName(v)
		if err != nil {
			return err
		}
		i.WishName = wn
		return nil
	default:
		return errors.New("received invalid type for wish name")
	}

}

func (i WishName) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.WishName.String())
}
