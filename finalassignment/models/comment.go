package models

type Comment struct {
	Model
	UserID  uint   `json:"user_id" example:"1"`
	PhotoID uint   `json:"photo_id" example:"1"`
	Message string `gorm:"not null;type:varchar(8192)" json:"message"`
}
