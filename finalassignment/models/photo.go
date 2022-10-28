package models

type Photo struct {
	Model
	Title    string `gorm:"not null" json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `gorm:"not null" json:"photo_url" example:"https://subdomain.domain.dom.ge/path?arg=1"`
	UserID   uint   `json:"user_id" example:"1"`
}
