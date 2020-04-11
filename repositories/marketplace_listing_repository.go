package repositories

import (
	"encoding/base64"
	"errors"
	"github.com/jinzhu/gorm"
	m "github.com/214alphadev/marketplace-bl/entities"
	"github.com/214alphadev/marketplace-bl/repository"
	"github.com/214alphadev/marketplace-bl/value_objects"
)

type marketplaceListingModel struct {
	InternalID    uint `gorm:"AUTO_INCREMENT"`
	Seller        string
	ID            string `gorm:"unique,primary_key"`
	Name          string
	Description   string
	CreatedAt     int64
	PriceAmount   uint64
	PriceCurrency string
	Photo         string
}

func toMarketplaceListingEntity(listing marketplaceListingModel) (*m.Listing, error) {

	seller, err := value_objects.NewMemberID(listing.Seller)
	if err != nil {
		return nil, err
	}

	listingID, err := m.NewListingID(listing.ID)
	if err != nil {
		return nil, err
	}

	listingName, err := m.NewListingName(listing.Name)
	if err != nil {
		return nil, err
	}

	description, err := m.NewListingDescription(listing.Description)
	if err != nil {
		return nil, err
	}

	createdAt, err := m.NewListingCreationDate(listing.CreatedAt)
	if err != nil {
		return nil, err
	}

	price, err := m.NewPrice(listing.PriceCurrency, listing.PriceAmount)
	if err != nil {
		return nil, err
	}

	listingPhoto, err := base64.StdEncoding.DecodeString(listing.Photo)
	if err != nil {
		return nil, err
	}

	l, err := m.NewListing(seller, listingID, listingName, description, createdAt, price, listingPhoto)
	if err != nil {
		return &l, err
	}

	return &l, nil

}

func (vm marketplaceListingModel) TableName() string {
	return "marketplace_listing"
}

type marketplaceListingRepository struct {
	db *gorm.DB
}

func (r marketplaceListingRepository) Save(listing m.Listing) error {

	return r.db.Save(&marketplaceListingModel{
		Seller:        listing.Seller().String(),
		ID:            listing.ID().String(),
		Name:          listing.Name().String(),
		Description:   listing.Description().String(),
		CreatedAt:     listing.CreatedAt().Time().Unix(),
		PriceAmount:   listing.Price().Amount(),
		PriceCurrency: listing.Price().Currency(),
		Photo:         base64.StdEncoding.EncodeToString(listing.Photo()),
	}).Error
}

func (r marketplaceListingRepository) Get(listingID m.ListingID) (*m.Listing, error) {

	listing := &marketplaceListingModel{}

	err := r.db.Find(listing, "id = ?", listingID.String()).Error

	switch err {
	case nil:
		return toMarketplaceListingEntity(*listing)
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		return nil, err
	}

}

func (r marketplaceListingRepository) Query(from *m.ListingID, next uint32) ([]m.Listing, error) {

	fetchedListings := []marketplaceListingModel{}
	convertFetchedListings := func(listings []marketplaceListingModel) ([]m.Listing, error) {

		convertedFetchedListings := []m.Listing{}

		for _, w := range listings {
			mappedListings, err := toMarketplaceListingEntity(w)
			if err != nil {
				return nil, err
			}
			convertedFetchedListings = append(convertedFetchedListings, *mappedListings)
		}

		return convertedFetchedListings, nil

	}

	switch from {
	case nil:
		if err := r.db.Limit(next).Find(&fetchedListings).Error; err != nil {
			return nil, err
		}
		return convertFetchedListings(fetchedListings)
	default:

		listing := &marketplaceListingModel{}
		if err := r.db.Find(&listing, "id = ?", from.String()).Error; err != nil {
			return nil, err
		}

		if err := r.db.Find(&fetchedListings, "internal_id > ?", listing.ID).Error; err != nil {
			return nil, err
		}

		if err := r.db.Where("internal_id > ?", listing.InternalID).Limit(next).Find(&fetchedListings).Error; err != nil {
			return nil, err
		}

		return convertFetchedListings(fetchedListings)

	}

}

func (r marketplaceListingRepository) Delete(listingID m.ListingID) error {
	return errors.New("delete is not implemented")
}

func NewMarketplaceListingRepo(db *gorm.DB) (repository.IListingRepository, error) {

	if err := db.AutoMigrate(&marketplaceListingModel{}).Error; err != nil {
		return nil, err
	}

	return marketplaceListingRepository{
		db: db,
	}, nil

}
