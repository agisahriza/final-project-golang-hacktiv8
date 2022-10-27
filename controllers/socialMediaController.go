package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         		SocialMedia.ID,
		"name":      		SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":    		SocialMedia.UserID,
		"created_at": 		SocialMedia.CreatedAt,
	})
}

func GetAllSocialMedias(c *gin.Context) {
	db := database.GetDB()

	var SocialMedias []models.SocialMedia
	PrintSocialMedias := make([]models.PrintSocialMedia, 0)

	err := db.Debug().Find(&SocialMedias).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for _, socialMedia := range SocialMedias {
		var PrintSocialMedia models.PrintSocialMedia
		var User models.User
		var Photo models.Photo
		_ = db.First(&User, "id = ?", socialMedia.UserID).Error
		_ = db.First(&Photo, "id = ?", socialMedia.UserID).Error
		PrintSocialMedia.ID = socialMedia.ID
		PrintSocialMedia.Name = socialMedia.Name
		PrintSocialMedia.SocialMediaUrl = socialMedia.SocialMediaUrl
		PrintSocialMedia.UserID = socialMedia.UserID
		PrintSocialMedia.CreatedAt = socialMedia.CreatedAt
		PrintSocialMedia.UpdatedAt = socialMedia.UpdatedAt
		PrintSocialMedia.User.ID = socialMedia.UserID
		PrintSocialMedia.User.Username = User.Username

		PrintSocialMedias = append(PrintSocialMedias, PrintSocialMedia)
	}
	
	c.JSON(http.StatusOK, PrintSocialMedias)
}

func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID
	SocialMedia.ID = uint(socialMediaId)

	err := db.Debug().Model(&SocialMedia).Where("id = ?", socialMediaId).Updates(models.SocialMedia{Name: SocialMedia.Name, SocialMediaUrl: SocialMedia.SocialMediaUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SocialMedia)
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	
	SocialMedia := models.SocialMedia{}
	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))

	err := db.Where("id = ?", socialMediaId).Delete(&SocialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}