package marketplace

import (
	"encoding/json"
	"errors"
)

type Username struct {
	username string
}

func (Username) ImplementsGraphQLType(name string) bool {
	return name == "Username"
}

func (e *Username) UnmarshalGraphQL(input interface{}) error {
	return errors.New("can not unmarshal")
}

func (e Username) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.username)
}
