package marketplace

import (
	"encoding/json"
	"errors"
	"github.com/214alphadev/marketplace-bl/entities"
)

type ListingID struct {
	id entities.ListingID
}

func (ListingID) ImplementsGraphQLType(name string) bool {
	return name == "ListingID"
}

func (i *ListingID) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case string:
		listingID, err := entities.NewListingID(v)
		if err != nil {
			return err
		}
		i.id = listingID
		return nil
	default:
		return errors.New("failed to unmarshal listing id")
	}

}

func (i ListingID) String() string {
	return i.id.String()
}

func (i ListingID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.id.String())
}
