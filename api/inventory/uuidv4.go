package inventory

import (
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
)

type UUIDV4 struct {
	uuid.UUID
}

func (UUIDV4) ImplementsGraphQLType(name string) bool {
	return name == "UUIDV4"
}

func (i *UUIDV4) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case *string:
		id, err := uuid.FromString(*v)
		if err != nil {
			return err
		}
		i.UUID = id
		return nil
	case string:
		id, err := uuid.FromString(v)
		if err != nil {
			return err
		}
		i.UUID = id
		return nil
	default:
		return errors.New("failed to unmarshal uuid v4")
	}

}

func (i UUIDV4) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.UUID.String())
}
