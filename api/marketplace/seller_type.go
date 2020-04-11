package marketplace

import (
	cd "github.com/214alphadev/community-bl"
)

type Seller struct {
	member cd.MemberEntity
}

func (s Seller) EmailAddress() EmailAddress {
	return EmailAddress{
		emailAddress: s.member.EmailAddress.String(),
	}
}

func (s Seller) Username() Username {
	return Username{
		username: s.member.Username.String(),
	}
}

func (s Seller) FirstName() string {
	return s.member.Metadata.ProperName.FirstName()
}

func (s Seller) LastName() string {
	return s.member.Metadata.ProperName.LastName()
}

func (s Seller) ID() UUIDV4 {
	return UUIDV4{UUID: s.member.ID}
}

func (s Seller) ProfilePicture() *SellerProfilePicture {

	switch s.member.Metadata.ProfileImage {
	case nil:
		return nil
	default:
		return &SellerProfilePicture{
			picture: s.member.Metadata.ProfileImage.String(),
		}
	}

}

func newSeller(member cd.MemberEntity) Seller {
	return Seller{
		member: member,
	}
}
