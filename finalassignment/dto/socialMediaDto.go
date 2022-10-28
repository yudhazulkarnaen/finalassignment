package dto

type SocialMedia struct {
	Name           string `validate:"required"`
	SocialMediaUrl string `validate:"required" json:"social_media_url" example:"https://subdomain.domain.dom.ge/path"`
}
