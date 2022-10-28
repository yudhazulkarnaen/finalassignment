package database

import (
	"time"

	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
	"finalassignment.id/finalassignment/utils/token"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordMismatch = bcrypt.ErrMismatchedHashAndPassword

func GenerateToken(userDto dto.UserLogin) (jwt string, err error) {
	user, err := getUserByEmail(userDto.Email)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password))
	if err != nil {
		return
	}
	jwt, err = token.GenerateToken(user.ID)
	return
}
func DeleteUserById(id uint) error {
	user, err := GetUserWithoutPreload(id)
	if err != nil {
		return err
	}
	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}
func getUserByEmail(email string) (models.User, error) {
	user := models.User{}
	if db == nil {
		return user, ErrDbNotStarted
	}
	err := db.Model(&models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func GetUsernameAndEmail(id uint) (dto.UserUpdate, error) {
	userDto := dto.UserUpdate{}
	if db == nil {
		return userDto, ErrDbNotStarted
	}
	user := models.User{}
	if err := db.Select("username", "email").Take(&user, id).Error; err != nil {
		return userDto, err
	}
	userDto.Username = user.Username
	userDto.Email = user.Email
	return userDto, nil
}
func CreateUser(userRegister *dto.UserRegister) (ID uint, err error) {
	if db == nil {
		err = ErrDbNotStarted
		return
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 4)
	if err != nil {
		return
	}
	newUser := models.User{
		Username: userRegister.Username,
		Email:    userRegister.Email,
		Password: string(passwordBytes),
		Age:      userRegister.Age,
		Model: models.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = db.Create(&newUser).Error
	if err != nil {
		return
	}
	ID = newUser.ID
	return
}
func UpdateUser(id uint, userDto *dto.UserUpdate) (models.User, error) {
	user, err := GetUserWithoutPreload(id)
	if err != nil {
		return user, err
	}
	if userDto.Email != "" {
		user.Email = userDto.Email
	}
	if userDto.Username != "" {
		user.Username = userDto.Username
	}
	user.UpdatedAt = time.Now()
	err = db.Save(&user).Error
	return user, err
}
func GetUserWithoutPreload(id uint) (models.User, error) {
	user := models.User{}
	if db == nil {
		return user, ErrDbNotStarted
	}
	err := db.Model(&models.User{}).Take(&user, id).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
