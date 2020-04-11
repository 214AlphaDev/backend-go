package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	. "github.com/214alphadev/wishlist-bl/entities"
	. "github.com/214alphadev/wishlist-bl/repository"
	"github.com/214alphadev/wishlist-bl/value_objects"
)

type wishModel struct {
	ID          string `gorm:"unique,primary_key"`
	InternalID  uint   `gorm:"AUTO_INCREMENT"`
	Name        string
	Description string
	Story       *string
	Creator     string
	Category    string
}

func (m wishModel) TableName() string {
	return "wish_wishes"
}

func mapWishModelToEntity(wishModel wishModel) (*Wish, error) {

	memberID, err := value_objects.NewMemberID(wishModel.Creator)
	if err != nil {
		return nil, err
	}

	wishID, err := NewWishID(wishModel.ID)
	if err != nil {
		return nil, err
	}

	wishName, err := NewWishName(wishModel.Name)
	if err != nil {
		return nil, err
	}

	wishDescription, err := NewWishDescription(wishModel.Description)
	if err != nil {
		return nil, err
	}

	wishStory, err := NewWishStory(wishModel.Story)
	if err != nil {
		return nil, err
	}

	var category value_objects.Category
	switch wishModel.Category {
	case "Book":
		category = value_objects.Book
	case "Seed":
		category = value_objects.Seed
	case "Other":
		category = value_objects.Other
	}

	w, err := NewWish(memberID, wishID, wishName, wishDescription, wishStory, category)
	if err != nil {
		return nil, err
	}

	return &w, nil

}

type wishRepository struct {
	db *gorm.DB
}

func (wr *wishRepository) Save(wish Wish) error {

	w := &wishModel{
		ID:          wish.ID().String(),
		Name:        wish.Name().String(),
		Description: wish.Description().String(),
		Creator:     wish.Creator().String(),
	}

	switch wish.Category() {
	case value_objects.Seed:
		w.Category = "Seed"
	case value_objects.Book:
		w.Category = "Book"
	case value_objects.Other:
		w.Category = "Other"
	default:
		return fmt.Errorf("couldn't map category: %s", wish.Category())
	}

	if !wish.Story().IsNil() {
		s := wish.Story().String()
		w.Story = &s
	}

	return wr.db.Save(w).Error

}

func (wr *wishRepository) Get(wishID WishID) (*Wish, error) {

	wish := &wishModel{}

	err := wr.db.Find(wish, "id = ?", wishID.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapWishModelToEntity(*wish)
	default:
		return nil, err
	}

}

func (wr *wishRepository) Query(from *WishID, next uint32) ([]Wish, error) {

	fetchedWishes := []wishModel{}
	convertFetchedWishes := func(wishes []wishModel) ([]Wish, error) {

		convertedFetchedWishes := []Wish{}

		for _, w := range fetchedWishes {
			mappedWish, err := mapWishModelToEntity(w)
			if err != nil {
				return nil, err
			}
			convertedFetchedWishes = append(convertedFetchedWishes, *mappedWish)
		}

		return convertedFetchedWishes, nil

	}

	switch from {
	case nil:
		if err := wr.db.Limit(next).Find(&fetchedWishes).Error; err != nil {
			return nil, err
		}
		return convertFetchedWishes(fetchedWishes)
	default:

		item := &wishModel{}
		if err := wr.db.Find(&item, "id = ?", from.String()).Error; err != nil {
			return nil, err
		}

		if err := wr.db.Find(&fetchedWishes, "internal_id > ?", item.InternalID).Error; err != nil {
			return nil, err
		}

		if err := wr.db.Where("internal_id > ?", item.InternalID).Limit(next).Find(&fetchedWishes).Error; err != nil {
			return nil, err
		}

		return convertFetchedWishes(fetchedWishes)

	}

}

func (wr *wishRepository) VotesOf(wishID WishID) (uint32, error) {

	wishes := uint32(0)

	err := wr.db.Model(&voteModel{}).Where("wish_id = ?", wishID.String()).Count(&wishes).Error

	switch err {
	case nil:
		return wishes, nil
	default:
		return 0, err
	}

}

func NewWishRepository(db *gorm.DB) (IWishRepository, error) {

	if err := db.AutoMigrate(&wishModel{}).Error; err != nil {
		return nil, err
	}

	return &wishRepository{
		db: db,
	}, nil
}
