package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	GormModel
	Name           string `gorm:"not null" json:"name" form:"name" valid:"required~Your name is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url" valid:"required~Social media URL is required"`
	UserID         uint
}

type PrintSocialMedia struct {
	DefaultModel
	Name string `json:"name"`
	SocialMediaUrl string `json:"social_media_url"`
	UserID uint `json:"user_id"`
	User UserForSocialMedia
}

type UserForSocialMedia struct {
	ID uint `json:"id"`
	Username string `json:"username"`
}

func (s *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(s)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

