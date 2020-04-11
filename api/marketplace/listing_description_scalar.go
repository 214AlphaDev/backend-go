package marketplace

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/marketplace-bl/entities"
)

type ListingDescription struct {
	entities.ListingDescription
}

func (ListingDescription) ImplementsGraphQLType(name string) bool {
	return name == "ListingDescription"
}

func (i *ListingDescription) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		wd, err := entities.NewListingDescription(*v)
		if err != nil {
			return err
		}
		i.ListingDescription = wd
		return nil
	case string:
		wd, err := entities.NewListingDescription(v)
		if err != nil {
			return err
		}
		i.ListingDescription = wd
		return nil
	default:
		return errors.New("received invalid type for listing description")
	}

}

func (i ListingDescription) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ListingDescription.String())
}
