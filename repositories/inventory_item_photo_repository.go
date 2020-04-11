package repositories

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	. "github.com/214alphadev/inventory-bl/entities"
	. "github.com/214alphadev/inventory-bl/repository"
)

type itemPhotoModel struct {
	ID     string `gorm:"unique"`
	Photo  string
	ItemID string
}

func (m itemPhotoModel) TableName() string {
	return "inventory_item_photo"
}

type inventoryItemPhotoRepository struct {
	db *gorm.DB
}

func (r inventoryItemPhotoRepository) Save(photo ItemPhoto) error {
	return r.db.Save(&itemPhotoModel{
		ID:     photo.ID().String(),
		Photo:  photo.String(),
		ItemID: photo.Item().String(),
	}).Error
}

func (r inventoryItemPhotoRepository) DeleteAllFor(item ItemID) error {
	return r.db.Delete(&itemPhotoModel{}, "item_id = ?", item.String()).Error
}

func (r inventoryItemPhotoRepository) Get(item ItemID) (*ItemPhoto, error) {

	i := &itemPhotoModel{}
	err := r.db.Where("item_id = ?", item.String()).Find(&i).Error

	switch err {
	case nil:

		itemPhotoID, err := NewItemPhotoID(i.ID)
		if err != nil {
			return nil, err
		}

		itemID, err := NewItemID(i.ItemID)
		if err != nil {
			return nil, err
		}

		rawPhoto, err := base64.StdEncoding.DecodeString(i.Photo)
		if err != nil {
			return nil, err
		}

		itemPhoto, err := NewItemPhoto(itemPhotoID, itemID, rawPhoto)
		if err != nil {
			return nil, err
		}

		return &itemPhoto, nil

	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		return nil, err
	}

}

func NewInventoryItemPhotoRepository(db *gorm.DB) (IItemPhotoRepository, error) {

	if err := db.AutoMigrate(&itemPhotoModel{}).Error; err != nil {
		return nil, err
	}

	return &inventoryItemPhotoRepository{
		db: db,
	}, nil

}
