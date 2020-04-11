package marketplace

import "github.com/214alphadev/marketplace-bl/entities"

type Price struct {
	price entities.Price
}

func (p Price) Amount() PriceAmount {
	return PriceAmount{
		amount: p.price.Amount(),
	}
}

func (p Price) Currency() Currency {
	return Currency{
		currency: p.price.Currency(),
	}
}

func newPrice(p entities.Price) Price {
	return Price{
		price: p,
	}
}
