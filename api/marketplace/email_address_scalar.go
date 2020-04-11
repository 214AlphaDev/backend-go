package marketplace

import (
	"encoding/json"
	"errors"
)

type EmailAddress struct {
	emailAddress string
}

func (EmailAddress) ImplementsGraphQLType(name string) bool {
	return name == "EmailAddress"
}

func (e *EmailAddress) UnmarshalGraphQL(input interface{}) error {
	return errors.New("can not unmarshal")
}

func (e EmailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.emailAddress)
}
