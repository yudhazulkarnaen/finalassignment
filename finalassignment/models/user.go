package models

type User struct {
	Model
	Username     string        `gorm:"not null;uniqueIndex"`
	Email        string        `gorm:"not null;uniqueIndex"`
	Password     string        `gorm:"not null"`
	Age          uint          `gorm:"not null"`
	Photos       []Photo       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Comments     []Comment     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SocialMedias []SocialMedia `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
