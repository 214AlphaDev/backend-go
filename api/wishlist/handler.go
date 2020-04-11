package wishlist

import (
	gql "github.com/graph-gophers/graphql-go"
	gqlh "github.com/214alphadev/graphql-handler"
	. "github.com/214alphadev/wishlist-bl/services"
)

const schema = `
schema {
    query: Query
    mutation: Mutation
}

scalar UUIDV4
scalar WishName
scalar WishDescription
scalar WishStory

type Query {
    # get all wishes
    wishes(start: UUIDV4, next: Int!) : [Wish!]!
    # get the amount of votes left returns nil on error
    votesLeft() : Int
}

type Wish {
    # unique id of the wish
    id: UUIDV4!
    # name of the wish
    name: WishName!
    # description of the wish
    description: WishDescription!
    # story behind this wish
    story: WishStory
    # did the currently authenticated person voted on it
    votedOnIt: Boolean!
    # votes on that wish
    votes: Int!
    # category of wish
    category: Category!
}

enum CreateWishError {
    Unauthenticated
}

type CreateWishResponse {
    # the created wish
    wish: Wish
    # error that happend during wish creation
    error: CreateWishError
}

enum VoteError {
    WishDoesNotExist
    NoVotesLeft
    Unauthenticated
}

type VoteResponse {
    # error that might happen during voting
    error: VoteError
    # votes that are left
    votesLeft: Int
    # wish
    wish: Wish
}

enum WithdrawVoteError {
    WishDoesNotExist
    Unauthenticated
}

enum Category {
    Book
    Seed
    Other
}

type WithdrawVoteResponse {
    # error happend during withdrawing a vote
    error: WithdrawVoteError
    # votes that the current logged in user has left
    votesLeft: Int
    # wish
    wish: Wish
}

type Mutation {
    # create a new wish
    create(name: WishName!, description: WishDescription!, story: WishStory, category: Category!) : CreateWishResponse!
    # vote on a wish
    vote(id: UUIDV4!) : VoteResponse!
    # remove vote from wish
    withdrawVote(id: UUIDV4!) : WithdrawVoteResponse!
}
`

func NewGraphqlHandler(wishService IWishService, voteService IVoteService, logger gqlh.Logger) (*gqlh.Handler, error) {

	resolver := newResolver(wishService, voteService)

	schema, err := gql.ParseSchema(schema, resolver)
	if err != nil {
		return nil, err
	}

	return gqlh.NewHandler(schema, logger)

}
