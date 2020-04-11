package marketplace

import (
	"encoding/json"
	"errors"
)

type SellerProfilePicture struct {
	picture string
}

func (SellerProfilePicture) ImplementsGraphQLType(name string) bool {
	return name == "SellerProfilePicture"
}

func (s *SellerProfilePicture) UnmarshalGraphQL(input interface{}) error {
	return errors.New("can not unmarshal")
}

func (s SellerProfilePicture) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.picture)
}
