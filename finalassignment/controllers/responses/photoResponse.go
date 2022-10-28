package responses

import (
	"time"

	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
)

type Photo struct {
	ID       uint   `json:"id" example:"1"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" example:"https://subdomain.domain.dom.ge/path?arg=1"`
	UserID   uint   `json:"user_id" example:"1"`
}

type CreatePhoto struct {
	Photo
	CreatedAt time.Time `json:"created_at" example:"2019-11-09T21:21:46+00:00"`
}

type GetPhoto struct {
	models.Photo
	User dto.UserUpdate
}

type UpdatePhoto struct {
	Photo
	UpdatedAt time.Time `json:"updated_at" example:"2019-11-09T21:21:46+00:00"`
}

func (getPhoto *GetPhoto) Set(photo models.Photo) {
	getPhoto.ID = photo.ID
	getPhoto.CreatedAt = photo.CreatedAt
	getPhoto.UpdatedAt = photo.UpdatedAt
	getPhoto.Title = photo.Title
	getPhoto.Caption = photo.Caption
	getPhoto.PhotoUrl = photo.PhotoUrl
	getPhoto.UserID = photo.UserID
}
