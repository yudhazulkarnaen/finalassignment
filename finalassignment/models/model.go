package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"PrimaryKey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2019-11-09T21:21:46+00:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2019-11-09T21:21:46+00:00"`
}
