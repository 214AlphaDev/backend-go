package wishlist

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/wishlist-bl/entities"
)

type WishStory struct {
	entities.WishStory
}

func (WishStory) ImplementsGraphQLType(name string) bool {
	return name == "WishStory"
}

func (i *WishStory) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		ws, err := entities.NewWishStory(v)
		if err != nil {
			return err
		}
		i.WishStory = ws
		return nil
	case string:
		ws, err := entities.NewWishStory(&v)
		if err != nil {
			return err
		}
		i.WishStory = ws
		return nil
	default:
		return errors.New("received invalid type for wish story")
	}

}

func (i WishStory) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.WishStory.String())
}
