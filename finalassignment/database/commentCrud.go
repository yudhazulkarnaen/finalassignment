package database

import (
	"time"

	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
)

func GetAllComments() ([]models.Comment, error) {
	comments := make([]models.Comment, 1)
	if err := db.Model(&models.Comment{}).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
func DeleteComment(commentID, userID uint) error {
	comment, err := GetSingleComment(commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return ErrIllegalUpdate
	}
	if err := db.Delete(&comment, commentID).Error; err != nil {
		return err
	}
	return nil
}
func GetSingleComment(commentID uint) (models.Comment, error) {
	comment := models.Comment{}
	if db == nil {
		return comment, ErrDbNotStarted
	}
	err := db.Model(&models.Comment{}).Take(&comment, commentID).Error
	return comment, err
}
func UpdateComment(commentID, userID uint, messageDto *dto.CommentMessage) (comment models.Comment, err error) {
	comment, err = GetSingleComment(commentID)
	if err != nil {
		return
	}
	if comment.UserID != userID {
		comment = models.Comment{}
		err = ErrIllegalUpdate
		return
	}
	comment.Message = messageDto.Message
	comment.UpdatedAt = time.Now()
	err = db.Save(&comment).Error
	return
}
func CreateComment(userID uint, commentDto *dto.Comment) (models.Comment, error) {
	if db == nil {
		return models.Comment{}, ErrDbNotStarted
	}
	newComment := models.Comment{
		UserID:  userID,
		PhotoID: commentDto.PhotoID,
		Message: commentDto.Message,
		Model: models.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	if err := db.Create(&newComment).Error; err != nil {
		return models.Comment{}, err
	}
	return newComment, nil
}
