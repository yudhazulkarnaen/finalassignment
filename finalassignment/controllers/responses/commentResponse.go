package responses

import (
	"time"

	"finalassignment.id/finalassignment/models"
)

type CreateComment struct {
	ID        uint      `json:"id" example:"1"`
	Message   string    `json:"message"`
	PhotoID   uint      `json:"photo_id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2019-11-09T21:21:46+00:00"`
}

type GetComment struct {
	CreateComment
	UpdatedAt time.Time `json:"updated_at" example:"2019-11-09T21:21:46+00:00"`
	User      UserComment
	Photo     models.Photo
}

type UserComment struct {
	ID       uint   `json:"id" example:"1"`
	Email    string `json:"email" example:"name@org.dom.ge"`
	Username string `json:"username"`
}

func (getComment *GetComment) Set(comment models.Comment) {
	getComment.ID = comment.ID
	getComment.CreatedAt = comment.CreatedAt
	getComment.UpdatedAt = comment.UpdatedAt
	getComment.UserID = comment.UserID
	getComment.PhotoID = comment.PhotoID
	getComment.Message = comment.Message
}
