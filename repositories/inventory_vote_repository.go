package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/214alphadev/inventory-bl/entities"
	"github.com/214alphadev/inventory-bl/repository"
)

type inventoryVoteModel struct {
	ItemID   string
	MemberID string
}

func (vm inventoryVoteModel) TableName() string {
	return "inventory_votes"
}

type inventoryVoteRepository struct {
	db *gorm.DB
}

func (r *inventoryVoteRepository) Save(v entities.Vote) error {
	return r.db.Save(&inventoryVoteModel{
		ItemID:   v.ItemID().String(),
		MemberID: v.MemberID().String(),
	}).Error
}

func (r *inventoryVoteRepository) DoesExist(v entities.Vote) (bool, error) {

	vote := &inventoryVoteModel{}

	err := r.db.Find(vote, "item_id = ? AND member_id = ?", v.ItemID().String(), v.MemberID().String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}

}

func (r *inventoryVoteRepository) Delete(v entities.Vote) error {
	return r.db.Delete(&inventoryVoteModel{}, "item_id = ? AND member_id = ?", v.ItemID().String(), v.MemberID().String()).Error
}

func NewInventoryVoteRepository(db *gorm.DB) (repository.IVoteRepository, error) {

	if err := db.AutoMigrate(&inventoryVoteModel{}).Error; err != nil {
		return nil, err
	}

	return &inventoryVoteRepository{
		db: db,
	}, nil

}
