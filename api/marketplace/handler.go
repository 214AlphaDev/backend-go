package marketplace

import (
	gql "github.com/graph-gophers/graphql-go"
	cd "github.com/214alphadev/community-bl"
	gqlh "github.com/214alphadev/graphql-handler"
	"github.com/214alphadev/marketplace-bl/services"
)

const schema = `
schema {
    query: Query
    mutation: Mutation
}

scalar UUIDV4
scalar ListingName
scalar ListingDescription
scalar Base64ListingPhoto
scalar ListingID
scalar PriceAmount
scalar EmailAddress
scalar Time
scalar SellerProfilePicture
scalar Username

type Query {
    # get all listings
    listings(start: ListingID, next: Int!) : [Listing!]!
}

enum Currency {
    USD
}

type Price {
    currency: Currency!
    amount: PriceAmount!
}

type Listing {
    id: ListingID!
    name: ListingName!
    description: ListingDescription!
    photo: Base64ListingPhoto!
    price: Price!
    listedAt: Time!
    seller: Seller!
}

enum ListNewListingError {
    Unauthenticated
}

type ListNewListingResponse {
    # the created listing
    listing: Listing
    # error that happend during listing listing
    error: ListNewListingError
}

input UpdateListingInput {
    id: ListingID!
    name: ListingName!
    description: ListingDescription!
    photo: Base64ListingPhoto!
    priceCurrency: Currency!
    priceAmount: PriceAmount!
}

enum UpdateListingError {
    ListingDoesNotExist
    Unauthorized
    Unauthenticated
}

type Seller {
    id: UUIDV4!
    emailAddress: EmailAddress!
    username: Username!
    firstName: String!
    lastName: String!
    profilePicture: SellerProfilePicture
}

type UpdateListingResponse {
    error: UpdateListingError
    listing: Listing
}

type Mutation {
    # list new listing
    listNewListing(name: ListingName!, description: ListingDescription!, photo: Base64ListingPhoto!, priceCurrency: Currency!, priceAmount: PriceAmount!) : ListNewListingResponse!
    # update a listing
    updateListing(input: UpdateListingInput!) : UpdateListingResponse!
}
`

func NewGraphqlHandler(listingService services.IListingService, communityService cd.CommunityInterface, logger gqlh.Logger) (*gqlh.Handler, error) {

	resolver := newResolver(listingService, communityService)

	schema, err := gql.ParseSchema(schema, resolver)
	if err != nil {
		return nil, err
	}

	return gqlh.NewHandler(schema, logger)

}
