package marketplace

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/marketplace-bl/entities"
)

type ListingName struct {
	entities.ListingName
}

func (ListingName) ImplementsGraphQLType(name string) bool {
	return name == "ListingName"
}

func (i *ListingName) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wn, err := entities.NewListingName(*v)
		if err != nil {
			return err
		}
		i.ListingName = wn
		return nil
	case string:
		wn, err := entities.NewListingName(v)
		if err != nil {
			return err
		}
		i.ListingName = wn
		return nil
	default:
		return errors.New("received invalid type for listing name")
	}

}

func (i ListingName) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ListingName.String())
}
