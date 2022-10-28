package models

type SocialMedia struct {
	Model
	Name           string `gorm:"not null;type:varchar(8192)" json:"name"`
	SocialMediaUrl string `gorm:"not null;type:varchar(8192)" json:"social_media_url"`
	UserID         uint   `json:"user_id"`
}
