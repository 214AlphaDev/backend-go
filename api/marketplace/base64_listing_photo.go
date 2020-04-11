package marketplace

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Base64ListingPhoto struct {
	photo []byte
}

func (Base64ListingPhoto) ImplementsGraphQLType(name string) bool {
	return name == "Base64ListingPhoto"
}

func (i *Base64ListingPhoto) Bytes() []byte {
	return i.photo
}

func (i *Base64ListingPhoto) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		out, err := base64.StdEncoding.DecodeString(*v)
		if err != nil {
			return err
		}
		i.photo = out
		return nil
	case string:
		out, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		i.photo = out
		return nil
	default:
		return errors.New("invalid base64 string")
	}

}

func (i Base64ListingPhoto) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(i.photo))
}
