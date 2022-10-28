package database

import (
	"time"

	"finalassignment.id/finalassignment/dto"
	_ "finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
)

func UpdatePhoto(photoID, userID uint, photoDto *dto.Photo) (UpdatedAt time.Time, err error) {
	photo, err := GetSinglePhoto(photoID)
	if err != nil {
		return
	}
	if photo.UserID != userID {
		err = ErrIllegalUpdate
		return
	}
	if photoDto.Title != "" {
		photo.Title = photoDto.Title
	}
	if photoDto.PhotoUrl != "" {
		photo.PhotoUrl = photoDto.PhotoUrl
	}
	photo.Caption = photoDto.Caption
	photo.UpdatedAt = time.Now()
	err = db.Save(&photo).Error
	return
}
func DeletePhoto(photoID, userID uint) error {
	photo, err := GetSinglePhoto(photoID)
	if err != nil {
		return err
	}
	if photo.UserID != userID {
		return ErrIllegalUpdate
	}
	if err := db.Delete(&photo, photoID).Error; err != nil {
		return err
	}
	return nil
}
func CreatePhoto(userID uint, photoDto *dto.Photo) (ID uint, err error) {
	if db == nil {
		err = ErrDbNotStarted
		return
	}
	newPhoto := models.Photo{
		Title:    photoDto.Title,
		Caption:  photoDto.Caption,
		PhotoUrl: photoDto.PhotoUrl,
		UserID:   userID,
		Model: models.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = db.Create(&newPhoto).Error
	if err != nil {
		return
	}
	ID = newPhoto.ID
	return
}
func GetAllPhotos() ([]models.Photo, error) {
	photos := make([]models.Photo, 1)
	if err := db.Model(&models.Photo{}).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}
func GetSinglePhoto(photoID uint) (models.Photo, error) {
	photo := models.Photo{}
	if db == nil {
		return photo, ErrDbNotStarted
	}
	err := db.Model(&models.Photo{}).Take(&photo, photoID).Error
	return photo, err
}
