package database

import (
	"errors"
	"fmt"

	"finalassignment.id/finalassignment/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host             = "localhost"
	dbUser           = "postgres"
	dbPort           = "5432"
	dbName           = "finalassignmentdb"
	db               *gorm.DB
	err              error
	ErrDbNotStarted  error = errors.New("DB hasn't started yet.")
	ErrIllegalUpdate       = errors.New("This resource is not yours.")
)

func StartDB() {
	var password string
	fmt.Println("Enter db password (not hidden, be careful of shoulder surfing)")
	fmt.Scanln(&password)
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, dbUser, password, dbName, dbPort)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}
func GetDB() *gorm.DB {
	return db
}
