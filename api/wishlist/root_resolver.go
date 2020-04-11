package wishlist

import (
	. "github.com/214alphadev/wishlist-bl/services"
)

type Resolver struct {
	wishService IWishService
	voteService IVoteService
}

func newResolver(wishService IWishService, voteService IVoteService) *Resolver {
	return &Resolver{
		wishService: wishService,
		voteService: voteService,
	}
}
