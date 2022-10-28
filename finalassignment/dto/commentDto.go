package dto

type CommentMessage struct {
	Message string `validate:"required"`
}
type Comment struct {
	Message string `validate:"required"`
	PhotoID uint   `validate:"required" json:"photo_id" example:"1"`
}
