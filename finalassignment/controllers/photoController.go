package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"finalassignment.id/finalassignment/controllers/responses"
	"finalassignment.id/finalassignment/database"
	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/utils/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePhoto godoc
// @Summary      Create a Photo
// @Description  Create a Photo associated with the logged in user identified by bearer token.
// @Tags         photos
// @Accept       json
// @Produce      json
// @Param        user body dto.Photo true "JSON of the photo to be made. Caption is not mandatory."
// @Success      201  {object}  responses.CreatePhoto
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /photos [post]
// @Security	 BearerAuth
func CreatePhoto(ctx *gin.Context) {
	var newPhoto dto.Photo
	if err := ctx.ShouldBindJSON(&newPhoto); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&newPhoto); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ID, err := database.CreatePhoto(userID, &newPhoto)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, responses.CreatePhoto{
		Photo: responses.Photo{
			ID:       ID,
			Title:    newPhoto.Title,
			Caption:  newPhoto.Caption,
			PhotoUrl: newPhoto.PhotoUrl,
			UserID:   userID,
		},
		CreatedAt: time.Now(),
	})
}

// GetPhotos godoc
// @Summary      Get photos
// @Description  Get photos.
// @Tags         photos
// @Accept       json
// @Produce      json
// @Success      200  {object}  []responses.GetPhoto
// @Failure		 400 {object} responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /photos [get]
// @Security	 BearerAuth
func GetAllPhotos(ctx *gin.Context) {
	photos, err := database.GetAllPhotos()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	photosResponse := make([]responses.GetPhoto, len(photos))
	userDtos := make(map[uint]dto.UserUpdate)
	for i, photo := range photos {
		photosResponse[i].Set(photo)
		userDto, ok := userDtos[photo.UserID]
		if !ok {
			userDto, err = database.GetUsernameAndEmail(photo.UserID)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			userDtos[photo.UserID] = userDto
		}
		photosResponse[i].User = userDto
	}
	ctx.JSON(http.StatusOK, photosResponse)
}

// UpdatePhoto godoc
// @Summary      Update a photo
// @Description  Update a photo associated with logged in user.
// @Tags         photos
// @Accept       json
// @Produce      json
// @Param		 photoId path uint true "ID number of the photo"
// @Param        photo body dto.Photo true "New JSON of the photo."
// @Success      200  {object}  responses.UpdatePhoto
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /photos/{photoId} [put]
// @Security	 BearerAuth
func UpdatePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	parsedID, err := strconv.ParseUint(photoID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	var photoDto dto.Photo
	if err := ctx.ShouldBindJSON(&photoDto); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&photoDto); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	updatedAt, err := database.UpdatePhoto(uint(parsedID), userID, &photoDto)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Photo with ID %d is not found.", parsedID),
			})
			return
		}
		if errors.Is(err, database.ErrIllegalUpdate) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, responses.ErrorMessage{
				ErrorMessage: err.Error(),
			})
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, responses.UpdatePhoto{
		Photo: responses.Photo{
			ID:       uint(parsedID),
			Title:    photoDto.Title,
			Caption:  photoDto.Caption,
			PhotoUrl: photoDto.PhotoUrl,
			UserID:   userID,
		},
		UpdatedAt: updatedAt,
	})
}

// DeletePhoto godoc
// @Summary      Delete a photo
// @Description  Delete a photo associated with logged in user.
// @Tags         photos
// @Accept       json
// @Produce      json
// @Param		 photoId path uint true "ID number of the photo to be deleted"
// @Success      200  {object}  responses.Message
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /photos/{photoId} [delete]
// @Security	 BearerAuth
func DeletePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	parsedID, err := strconv.ParseUint(photoID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := database.DeletePhoto(uint(parsedID), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Photo with ID %d is not found.", parsedID),
			})
			return
		}
		if errors.Is(err, database.ErrIllegalUpdate) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, responses.ErrorMessage{
				ErrorMessage: err.Error(),
			})
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, responses.Message{
		Message: "Your photo has been successfully deleted",
	})
}
