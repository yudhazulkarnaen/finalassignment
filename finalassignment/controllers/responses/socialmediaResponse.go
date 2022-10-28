package responses

import (
	"time"

	"finalassignment.id/finalassignment/models"
)

type CreateSocialMedia struct {
	ID             uint      `json:"id" example:"1"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url" example:"https://subdomain.domain.dom.ge/path"`
	UserID         uint      `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetAllSocialMedias struct {
	SocialMedias []GetSocialMedia `json:"social_medias"`
}

type GetSocialMedia struct {
	models.SocialMedia
	User UserSocialMedia
}

type UpdateSocialMedia struct {
	ID             uint      `json:"id" example:"1"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url" example:"https://subdomain.domain.dom.ge/path"`
	UserID         uint      `json:"user_id" example:"1"`
	UpdatedAt      time.Time `json:"updated_at" example:"2019-11-09T21:21:46+00:00"`
}

type UserSocialMedia struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username"`
}

func (getSocmed *GetSocialMedia) Set(socmed models.SocialMedia) {
	getSocmed.ID = socmed.ID
	getSocmed.CreatedAt = socmed.CreatedAt
	getSocmed.UpdatedAt = socmed.UpdatedAt
	getSocmed.Name = socmed.Name
	getSocmed.SocialMediaUrl = socmed.SocialMediaUrl
	getSocmed.UserID = socmed.UserID
}
