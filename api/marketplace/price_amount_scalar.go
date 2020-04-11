package marketplace

import (
	"encoding/json"
	"errors"
	"strconv"
)

type PriceAmount struct {
	amount uint64
}

func (PriceAmount) ImplementsGraphQLType(name string) bool {
	return name == "PriceAmount"
}

func (c *PriceAmount) UnmarshalGraphQL(input interface{}) error {

	switch v := input.(type) {
	case string:

		priceAmount, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}

		c.amount = priceAmount

		return nil

	default:
		return errors.New("failed to unmarshal currency")
	}

}

func (c PriceAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(c.amount, 10))
}
