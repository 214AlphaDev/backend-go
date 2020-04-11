package inventory

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Base64ItemPhoto struct {
	photo []byte
}

func (Base64ItemPhoto) ImplementsGraphQLType(name string) bool {
	return name == "Base64ItemPhoto"
}

func (i *Base64ItemPhoto) Bytes() []byte {
	return i.photo
}

func (i *Base64ItemPhoto) UnmarshalGraphQL(input interface{}) error {

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

func (i Base64ItemPhoto) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(i.photo))
}
