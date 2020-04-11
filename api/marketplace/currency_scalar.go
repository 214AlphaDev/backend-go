package marketplace

import (
	"encoding/json"
	"errors"
)

type Currency struct {
	currency string
}

func (Currency) ImplementsGraphQLType(name string) bool {
	return name == "Currency"
}

func (c *Currency) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case string:

		switch v {
		case "USD":
			c.currency = v
			return nil
		default:
			return errors.New("failed to unmarshal currency - invalid currency")
		}

	default:
		return errors.New("failed to unmarshal currency")
	}

}

func (c Currency) String() string {
	return c.currency
}

func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.currency)
}
