package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"finalassignment.id/finalassignment/controllers/responses"
	"finalassignment.id/finalassignment/database"
	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/utils/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSocialMedia godoc
// @Summary      Create a social media
// @Description  Create a social media associated with the logged in user.
// @Tags         socialMedias
// @Accept       json
// @Produce      json
// @Param        socialMedia body dto.SocialMedia true "JSON of the social media to be made."
// @Success      201  {object}  responses.CreateSocialMedia
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /socialmedias [post]
// @Security	 BearerAuth
func CreateSocialMedia(ctx *gin.Context) {
	var newSocmed dto.SocialMedia
	if err := ctx.ShouldBindJSON(&newSocmed); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&newSocmed); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	socmed, err := database.CreateSocialMedia(userID, &newSocmed)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, responses.CreateSocialMedia{
		ID:             socmed.ID,
		Name:           socmed.Name,
		SocialMediaUrl: socmed.SocialMediaUrl,
		UserID:         userID,
		CreatedAt:      socmed.CreatedAt,
	})
}

// GetSocialMedias godoc
// @Summary      Get social medias
// @Description  Get social medias.
// @Tags         socialMedias
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.GetAllSocialMedias
// @Failure		 400 {object} responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /socialmedias [get]
// @Security	 BearerAuth
func GetAllSocialMedias(ctx *gin.Context) {
	socmeds, err := database.GetAllSocialMedias()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	socmedsResponse := make([]responses.GetSocialMedia, len(socmeds))
	userDtos := make(map[uint]dto.UserUpdate)
	for i, socmed := range socmeds {
		socmedsResponse[i].Set(socmed)
		userDto, ok := userDtos[socmed.UserID]
		if !ok {
			userDto, err = database.GetUsernameAndEmail(socmed.UserID)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			userDtos[socmed.UserID] = userDto
		}
		socmedsResponse[i].User = responses.UserSocialMedia{
			ID:       socmed.UserID,
			Username: userDto.Username,
		}
	}
	ctx.JSON(http.StatusOK, responses.GetAllSocialMedias{
		SocialMedias: socmedsResponse,
	})
}

// UpdateSocialMedia godoc
// @Summary      Update a social media
// @Description  Update a social media associated with logged in user.
// @Tags         socialMedias
// @Accept       json
// @Produce      json
// @Param		 socialMediaId path uint true "ID number of the social media"
// @Param        socialMedia body dto.SocialMedia true "New JSON of the social media."
// @Success      200  {object}  responses.UpdateSocialMedia
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /socialmedias/{socialMediaId} [put]
// @Security	 BearerAuth
func UpdateSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")
	parsedID, err := strconv.ParseUint(socialMediaID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	var socialMediaDto dto.SocialMedia
	if err := ctx.ShouldBindJSON(&socialMediaDto); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&socialMediaDto); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	updatedAt, err := database.UpdateSocialMedia(uint(parsedID), userID, &socialMediaDto)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Social media with ID %d is not found.", parsedID),
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
	ctx.JSON(http.StatusOK, responses.UpdateSocialMedia{
		ID:             uint(parsedID),
		Name:           socialMediaDto.Name,
		SocialMediaUrl: socialMediaDto.SocialMediaUrl,
		UserID:         userID,
		UpdatedAt:      updatedAt,
	})
}

// DeleteSocialMedia godoc
// @Summary      Delete a social media
// @Description  Delete a social media associated with logged in user.
// @Tags         socialMedias
// @Accept       json
// @Produce      json
// @Param		 socialMediaId path uint true "ID number of the social media to be deleted"
// @Success      200  {object}  responses.Message
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /socialmedias/{socialMediaId} [delete]
// @Security	 BearerAuth
func DeleteSocialMedia(ctx *gin.Context) {
	socmedID := ctx.Param("socialMediaId")
	parsedID, err := strconv.ParseUint(socmedID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := database.DeleteSocialMedia(uint(parsedID), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Social media with ID %d is not found.", parsedID),
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
		Message: "Your social media has been successfully deleted",
	})
}
