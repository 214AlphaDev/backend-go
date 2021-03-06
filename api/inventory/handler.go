package inventory

import (
	gql "github.com/graph-gophers/graphql-go"
	cd "github.com/214alphadev/community-bl"
	gqlh "github.com/214alphadev/graphql-handler"
	"github.com/214alphadev/inventory-bl/services"
)

const schema = `
schema {
    query: Query
    mutation: Mutation
}

scalar UUIDV4
scalar ItemName
scalar ItemDescription
scalar ItemStory
scalar Base64ItemPhoto

type Query {
    # get all items
    items(start: UUIDV4, next: Int!) : [Item!]!
}

type Item {
    # unique id of the item
    id: UUIDV4!
    # name of the item
    name: ItemName!
    # description of the item
    description: ItemDescription!
    # story behind this item
    story: ItemStory
    # did the currently authenticated person voted on it
    votedOnIt: Boolean!
    # votes on that item
    votes: Int!
    # photo of item
    photo: Base64ItemPhoto
    # category of item
    category: Category!
}

enum CreateItemError {
    Unauthenticated
    Unauthorized
}

type CreateItemResponse {
    # the created item
    item: Item
    # error that happend during item creation
    error: CreateItemError
}

enum VoteError {
    ItemDoesNotExist
    Unauthenticated
}

type VoteResponse {
    # error that might happen during voting
    error: VoteError
    # item
    item: Item
}

enum WithdrawVoteError {
    ItemDoesNotExist
    Unauthenticated
}

type WithdrawVoteResponse {
    # error happend during withdrawing a vote
    error: WithdrawVoteError
    # item
    item: Item
}

enum DeleteItemError {
    ItemDoesNotExist
    Unauthorized
    Unauthenticated
}

input UpdateItemInput {
    item: UUIDV4!
    name: ItemName!
    description: ItemDescription!
    story: ItemStory
    photo: Base64ItemPhoto
    category: Category!
}

enum UpdateItemError {
    ItemDoesNotExist
    Unauthorized
    Unauthenticated
}

type UpdateItemResponse {
    error: UpdateItemError
    item: Item
}

enum Category {
    Book
    Seed
    Other
}

type Mutation {
    # create a new item
    create(name: ItemName!, description: ItemDescription!, story: ItemStory, photo: Base64ItemPhoto, category: Category!) : CreateItemResponse!
    # vote on a item
    vote(id: UUIDV4!) : VoteResponse!
    # remove vote from item
    withdrawVote(id: UUIDV4!) : WithdrawVoteResponse!
    # delete item response
    deleteItem(id: UUIDV4!) : DeleteItemError
    # update item
    updateItem(input: UpdateItemInput!) : UpdateItemResponse!
}
`

func NewGraphqlHandler(itemService services.IItemService, communityService cd.CommunityInterface, logger gqlh.Logger) (*gqlh.Handler, error) {

	resolver := newResolver(itemService, communityService)

	schema, err := gql.ParseSchema(schema, resolver)
	if err != nil {
		return nil, err
	}

	return gqlh.NewHandler(schema, logger)

}
