package database

import (
	"time"

	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
)

func UpdateSocialMedia(socmedID, userID uint, socmedDto *dto.SocialMedia) (UpdatedAt time.Time, err error) {
	socmed, err := GetSingleSocialMedia(socmedID)
	if err != nil {
		return
	}
	if socmed.UserID != userID {
		err = ErrIllegalUpdate
		return
	}
	socmed.Name = socmedDto.Name
	socmed.SocialMediaUrl = socmedDto.SocialMediaUrl
	socmed.UpdatedAt = time.Now()
	err = db.Save(&socmed).Error
	return
}
func DeleteSocialMedia(socmedID, userID uint) error {
	socmed, err := GetSingleSocialMedia(socmedID)
	if err != nil {
		return err
	}
	if socmed.UserID != userID {
		return ErrIllegalUpdate
	}
	if err := db.Delete(&socmed, socmedID).Error; err != nil {
		return err
	}
	return nil
}
func CreateSocialMedia(userID uint, socmedDto *dto.SocialMedia) (models.SocialMedia, error) {
	if db == nil {
		return models.SocialMedia{}, ErrDbNotStarted
	}
	newSocmed := models.SocialMedia{
		UserID:         userID,
		Name:           socmedDto.Name,
		SocialMediaUrl: socmedDto.SocialMediaUrl,
		Model: models.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	if err := db.Create(&newSocmed).Error; err != nil {
		return models.SocialMedia{}, err
	}
	return newSocmed, nil
}
func GetAllSocialMedias() ([]models.SocialMedia, error) {
	socmeds := make([]models.SocialMedia, 1)
	if err := db.Model(&models.SocialMedia{}).Find(&socmeds).Error; err != nil {
		return nil, err
	}
	return socmeds, nil
}
func GetSingleSocialMedia(socmedID uint) (models.SocialMedia, error) {
	socmed := models.SocialMedia{}
	if db == nil {
		return socmed, ErrDbNotStarted
	}
	err := db.Model(&models.SocialMedia{}).Take(&socmed, socmedID).Error
	return socmed, err
}
