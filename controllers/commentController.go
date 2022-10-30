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

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID

	err := db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoID,
		"user_id":    Comment.UserID,
		"created_at": Comment.CreatedAt,
	})
}

func GetAllComments(c *gin.Context) {
	db := database.GetDB()

	var Comments []models.Comment
	PrintComments := make([]models.PrintComment, 0)

	err := db.Debug().Find(&Comments).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for _, comment := range Comments {
		var PrintComment models.PrintComment
		var User models.User
		var Photo models.Photo
		_ = db.First(&User, "id = ?", comment.UserID).Error
		_ = db.First(&Photo, "id = ?", comment.PhotoID).Error
		PrintComment.ID = comment.ID
		PrintComment.Message = comment.Message
		PrintComment.PhotoID = comment.PhotoID
		PrintComment.UserID = comment.UserID
		PrintComment.CreatedAt = comment.CreatedAt
		PrintComment.UpdatedAt = comment.UpdatedAt
		PrintComment.User.ID = User.ID
		PrintComment.User.Email = User.Email
		PrintComment.User.Username = User.Username
		PrintComment.Photo.ID = Photo.ID
		PrintComment.Photo.Title = Photo.Title
		PrintComment.Photo.Caption = Photo.Caption
		PrintComment.Photo.PhotoUrl = Photo.PhotoUrl
		PrintComment.Photo.UserID = Photo.UserID

		PrintComments = append(PrintComments, PrintComment)
	}
	
	c.JSON(http.StatusOK, PrintComments)
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.ID = uint(commentId)

	err := db.Debug().Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Comment)
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	
	Comment := models.Comment{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))

	err := db.Where("id = ?", commentId).Delete(&Comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}