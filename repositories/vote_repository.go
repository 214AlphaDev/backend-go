package repositories

import (
	"github.com/jinzhu/gorm"
	. "github.com/214alphadev/wishlist-bl/entities"
	. "github.com/214alphadev/wishlist-bl/repository"
	. "github.com/214alphadev/wishlist-bl/value_objects"
)

type voteModel struct {
	WishID   string
	MemberID string
}

func (vm voteModel) TableName() string {
	return "wish_votes"
}

type voteRepository struct {
	db *gorm.DB
}

func (r *voteRepository) Save(v Vote) error {
	return r.db.Save(&voteModel{
		WishID:   v.WishID().String(),
		MemberID: v.MemberID().String(),
	}).Error
}

func (r *voteRepository) DoesExist(v Vote) (bool, error) {

	vote := &voteModel{}

	err := r.db.Find(vote, "wish_id = ? AND member_id = ?", v.WishID().String(), v.MemberID().String()).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}

}

func (r *voteRepository) Delete(v Vote) error {
	return r.db.Delete(&voteModel{}, "wish_id = ? AND member_id = ?", v.WishID().String(), v.MemberID().String()).Error
}

func (r *voteRepository) GetVotesOfMember(memberID MemberID) ([]Vote, error) {

	votes := []voteModel{}

	if err := r.db.Find(&votes, "member_id = ?", memberID.String()).Error; err != nil {
		return nil, err
	}

	voteEntities := []Vote{}

	for _, vote := range votes {

		memberID, err := NewMemberID(vote.MemberID)
		if err != nil {
			return nil, err
		}

		wishID, err := NewWishID(vote.WishID)
		if err != nil {
			return nil, err
		}

		voteEntity, err := NewVote(memberID, wishID)
		if err != nil {
			return nil, err
		}

		voteEntities = append(voteEntities, voteEntity)

	}

	return voteEntities, nil

}

func NewVoteRepository(db *gorm.DB) (IVoteRepository, error) {

	if err := db.AutoMigrate(&voteModel{}).Error; err != nil {
		return nil, err
	}

	return &voteRepository{
		db: db,
	}, nil

}
