package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/justinas/alice"
	ca "github.com/214alphadev/community-api-go"
	am "github.com/214alphadev/community-authentication-middleware"
	cd "github.com/214alphadev/community-bl"
	cdi "github.com/214alphadev/community-infrastructure"
	"github.com/214alphadev/email-delivery-go"
	gqlh "github.com/214alphadev/graphql-handler"
	is "github.com/214alphadev/inventory-bl/services"
	ms "github.com/214alphadev/marketplace-bl/services"
	"github.com/214alphadev/wishlist-bl/services"
	"net/http"
	"s33d-backend/api/inventory"
	"s33d-backend/api/marketplace"
	"s33d-backend/api/wishlist"
	. "s33d-backend/config"
	"s33d-backend/repositories"
)

type ApplicationFactoryType = func(config Configuration) (*Application, error)

type consoleTransport struct{}

func (t *consoleTransport) SendConfirmationCode(confirmationCode cd.ConfirmationCode) error {
	fmt.Println("your confirmation code: ", confirmationCode.ConfirmationCode.String())
	return nil
}

func ApplicationFactory(config Configuration) (*Application, error) {

	db, err := gorm.Open("postgres", config.DatabaseURL())
	if err != nil {
		return nil, err
	}

	memberRepo, err := cdi.NewMemberRepository(db)
	if err != nil {
		return nil, err
	}

	applicationRepo, err := cdi.NewApplicationRepository(db)
	if err != nil {
		return nil, err
	}

	confCodeRepo, err := cdi.NewConfirmationCodeRepository(db)
	if err != nil {
		return nil, err
	}

	memberAccessPublicKeyRepo, err := cdi.NewMemberAccessPublicKeyRepository(db)
	if err != nil {
		return nil, err
	}

	accessTokenRepository, err := cdi.NewAccessTokenRepository(db)
	if err != nil {
		return nil, err
	}

	var communityTransport cd.Transport
	switch config.Development() {
	case true:
		communityTransport = &consoleTransport{}
	default:
		ss, err := ed.NewSendGridService(config.SendGridApiKey(), false)
		if err != nil {
			return nil, err
		}
		emailAddress, err := config.FromMail()
		if err != nil {
			return nil, err
		}
		communityTransport = cdi.NewTransport(ss, emailAddress)
	}

	community, err := cd.NewCommunity(cd.Dependencies{
		MemberRepository:                memberRepo,
		ApplicationRepository:           applicationRepo,
		ConfirmationCodeRepository:      confCodeRepo,
		Transport:                       communityTransport,
		MemberAccessPublicKeyRepository: memberAccessPublicKeyRepo,
		AccessTokenSigningKey:           config.AccessTokenSigningKey(),
		AccessTokenRepository:           accessTokenRepository,
	})
	if err != nil {
		return nil, err
	}

	authMiddleware := am.NewAuthenticateMemberMiddleware(community)

	itemRepo, err := repositories.NewItemRepository(db)
	if err != nil {
		return nil, err
	}
	voteRepo, err := repositories.NewInventoryVoteRepository(db)
	if err != nil {
		return nil, err
	}
	itemPhotoRepo, err := repositories.NewInventoryItemPhotoRepository(db)
	if err != nil {
		return nil, err
	}
	inventoryApiHandler, err := inventory.NewGraphqlHandler(
		is.NewItemService(itemRepo, voteRepo, itemPhotoRepo),
		community,
		gqlh.ConsoleLogger{},
	)
	if err != nil {
		return nil, err
	}

	communityApiHandler, err := ca.NewCommunityApi(community, gqlh.ConsoleLogger{})
	if err != nil {
		return nil, err
	}

	voteRepository, err := repositories.NewVoteRepository(db)
	if err != nil {
		return nil, err
	}

	wishRepository, err := repositories.NewWishRepository(db)
	if err != nil {
		return nil, err
	}

	voteService := services.NewVoteService(voteRepository)
	wishService := services.NewWishService(wishRepository, voteRepository, voteService)

	wishlistApiHandler, err := wishlist.NewGraphqlHandler(wishService, voteService, gqlh.ConsoleLogger{})
	if err != nil {
		return nil, err
	}

	listingRepo, err := repositories.NewMarketplaceListingRepo(db)
	if err != nil {
		return nil, err
	}
	handler, err := marketplace.NewGraphqlHandler(ms.NewListingService(listingRepo), community, gqlh.ConsoleLogger{})
	if err != nil {
		return nil, err
	}

	http.Handle("/api/community", alice.New(authMiddleware).Then(communityApiHandler))
	http.Handle("/api/inventory", alice.New(authMiddleware).Then(inventoryApiHandler))
	http.Handle("/api/wishlist", alice.New(authMiddleware).Then(wishlistApiHandler))
	http.Handle("/api/marketplace", alice.New(authMiddleware).Then(handler))

	return &Application{
		config:    config,
		Community: community,
	}, nil

}
