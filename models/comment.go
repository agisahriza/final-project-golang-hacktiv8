package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	GormModel
	UserID uint `json:"user_id"`
	PhotoID uint `json:"photo_id"`
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~Message is required"`
	User User
	Photo Photo
}

type PrintComment struct {
	DefaultModel
	UserID uint `json:"user_id"`
	PhotoID uint `json:"photo_id"`
	Message string `json:"message"`
	User UserForComment
	Photo PhotoForComment
}

type UserForComment struct {
	ID uint `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
}

type PhotoForComment struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Caption string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserID uint `json:"user_id"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
