package controllers

import (
	"errors"
	"net/http"

	"finalassignment.id/finalassignment/controllers/responses"
	"finalassignment.id/finalassignment/database"
	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Register a new user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body dto.UserRegister true "JSON of the user to be made. Minimum age is 9. Minimum password length is 6"
// @Success      201  {object}  responses.UserRegister
// @Success		 209  {object}  responses.Message
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /users/register [post]
func RegisterUser(ctx *gin.Context) {
	var newUser dto.UserRegister
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&newUser); err != nil {
		validationAbort(err, ctx)
		return
	}
	ID, err := database.CreateUser(&newUser)
	if err != nil {
		var perr *pgconn.PgError
		if ok := errors.As(err, &perr); ok {
			if perr.Code == uniqueViolationErr {
				ctx.AbortWithStatusJSON(209, responses.Message{
					Message: "The email or username is already registered. If it is yours, do login instead.",
				})
				return
			}
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, responses.UserRegister{
		Age:      newUser.Age,
		Email:    newUser.Email,
		ID:       ID,
		Username: newUser.Username,
	})
}

// LoginUser godoc
// @Summary      Login a user
// @Description  Login a user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body dto.UserLogin true "JSON of the user to login. Minimum password length is 6."
// @Success      200  {object}  responses.UserLogin
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /users/login [post]
func LoginUser(ctx *gin.Context) {
	var userLogin dto.UserLogin
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&userLogin); err != nil {
		validationAbort(err, ctx)
		return
	}
	jwt, err := database.GenerateToken(userLogin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, database.ErrPasswordMismatch) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.ErrorMessage{
				ErrorMessage: "Email or password is incorrect.",
			})
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, responses.UserLogin{
		Token: jwt,
	})
}

// UpdateUser godoc
// @Summary      Update logged in user
// @Description  update logged in user identified by their bearer token.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body dto.UserUpdate true "New email and new username of the logged in user. Leave one of it empty if you want it to stay the same."
// @Success      200  {object}  responses.UserUpdate
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /users [put]
// @Security	 BearerAuth
func UpdateUser(ctx *gin.Context) {
	var userDto dto.UserUpdate
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&userDto); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user, err := database.UpdateUser(userID, &userDto)
	if err != nil {
		var perr *pgconn.PgError
		if ok := errors.As(err, &perr); ok {
			if perr.Code == uniqueViolationErr {
				ctx.AbortWithStatusJSON(http.StatusOK, responses.Message{
					Message: "The email or username is already registered.",
				})
				return
			}
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, responses.UserUpdate{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Age:       user.Age,
		UpdatedAt: user.UpdatedAt,
	})
}

// DeleteOrder godoc
// @Summary      Delete logged in user
// @Description  Delete logged in user identified by their bearer token.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Message
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /users [delete]
// @Security	 BearerAuth
func DeleteUser(ctx *gin.Context) {
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := database.DeleteUserById(userID); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, responses.Message{
		Message: "Your account has been successfully deleted",
	})
}
