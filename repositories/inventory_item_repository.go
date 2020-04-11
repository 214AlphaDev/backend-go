package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	. "github.com/214alphadev/inventory-bl/entities"
	. "github.com/214alphadev/inventory-bl/repository"
	vo "github.com/214alphadev/inventory-bl/value_objects"
)

type inventoryItemModel struct {
	ID          string `gorm:"unique,primary_key"`
	InternalID  uint   `gorm:"AUTO_INCREMENT"`
	Name        string
	Description string
	Story       *string
	Creator     string
	CreatedAt   int64
	Category    string
}

func (m inventoryItemModel) TableName() string {
	return "inventory_items"
}

func mapItemModelToEntity(itemModel inventoryItemModel) (*Item, error) {

	memberID, err := vo.NewMemberID(itemModel.Creator)
	if err != nil {
		return nil, err
	}

	itemID, err := NewItemID(itemModel.ID)
	if err != nil {
		return nil, err
	}

	itemName, err := NewItemName(itemModel.Name)
	if err != nil {
		return nil, err
	}

	itemDescription, err := NewItemDescription(itemModel.Description)
	if err != nil {
		return nil, err
	}

	itemStory, err := NewItemStory(itemModel.Story)
	if err != nil {
		return nil, err
	}

	createdAt, err := NewItemCreationDate(itemModel.CreatedAt)
	if err != nil {
		return nil, err
	}

	var category vo.Category
	switch itemModel.Category {
	case "Seed":
		category = vo.Seed
	case "Book":
		category = vo.Book
	case "Other":
		category = vo.Other
	default:
		return nil, fmt.Errorf(`can't map category: "%s"'`, itemModel.Category)
	}

	w, err := NewItem(memberID, itemID, itemName, itemDescription, itemStory, createdAt, category)
	if err != nil {
		return nil, err
	}

	return &w, nil

}

type itemRepository struct {
	db *gorm.DB
}

func (wr *itemRepository) Save(item Item) error {

	w := &inventoryItemModel{
		ID:          item.ID().String(),
		Name:        item.Name().String(),
		Description: item.Description().String(),
		Creator:     item.Creator().String(),
		CreatedAt:   item.CreatedAt().Time().Unix(),
	}

	switch item.Category() {
	case vo.Book:
		w.Category = "Book"
	case vo.Seed:
		w.Category = "Seed"
	case vo.Other:
		w.Category = "Other"
	default:
		return fmt.Errorf("received invalid category: %s", item.Category())
	}

	if !item.Story().IsNil() {
		s := item.Story().String()
		w.Story = &s
	}

	return wr.db.Save(w).Error

}

func (wr *itemRepository) Get(itemID ItemID) (*Item, error) {

	item := &inventoryItemModel{}

	err := wr.db.Find(item, "id = ?", itemID.String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return mapItemModelToEntity(*item)
	default:
		return nil, err
	}

}

func (wr *itemRepository) Query(from *ItemID, next uint32) ([]Item, error) {

	fetchedItems := []inventoryItemModel{}
	convertFetchedItems := func(items []inventoryItemModel) ([]Item, error) {

		convertedFetchedItems := []Item{}

		for _, w := range items {
			mappedItems, err := mapItemModelToEntity(w)
			if err != nil {
				return nil, err
			}
			convertedFetchedItems = append(convertedFetchedItems, *mappedItems)
		}

		return convertedFetchedItems, nil

	}

	switch from {
	case nil:
		if err := wr.db.Limit(next).Find(&fetchedItems).Error; err != nil {
			return nil, err
		}
		return convertFetchedItems(fetchedItems)
	default:

		item := &inventoryItemModel{}
		if err := wr.db.Find(&item, "id = ?", from.String()).Error; err != nil {
			return nil, err
		}

		if err := wr.db.Find(&fetchedItems, "internal_id > ?", item.InternalID).Error; err != nil {
			return nil, err
		}

		if err := wr.db.Where("internal_id > ?", item.InternalID).Limit(next).Find(&fetchedItems).Error; err != nil {
			return nil, err
		}

		return convertFetchedItems(fetchedItems)

	}

}

func (wr *itemRepository) VotesOf(itemID ItemID) (uint32, error) {

	items := uint32(0)

	err := wr.db.Model(&inventoryVoteModel{}).Where("item_id = ?", itemID.String()).Count(&items).Error

	switch err {
	case nil:
		return items, nil
	default:
		return 0, err
	}

}

func (wr *itemRepository) Delete(itemID ItemID) error {
	return wr.db.Delete(&inventoryItemModel{}, "id = ?", itemID.String()).Error
}

func NewItemRepository(db *gorm.DB) (IItemRepository, error) {

	if err := db.AutoMigrate(&inventoryItemModel{}).Error; err != nil {
		return nil, err
	}

	return &itemRepository{
		db: db,
	}, nil
}
