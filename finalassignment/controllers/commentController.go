package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"finalassignment.id/finalassignment/controllers/responses"
	"finalassignment.id/finalassignment/database"
	"finalassignment.id/finalassignment/dto"
	"finalassignment.id/finalassignment/models"
	"finalassignment.id/finalassignment/utils/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment godoc
// @Summary      Create a Comment
// @Description  Create a Comment associated with the logged in user.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        comment body dto.Comment true "JSON of the comment to be made. Caption is not mandatory."
// @Success      201  {object}  responses.CreateComment
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /comments [post]
// @Security	 BearerAuth
func CreateComment(ctx *gin.Context) {
	var newComment dto.Comment
	if err := ctx.ShouldBindJSON(&newComment); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&newComment); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	comment, err := database.CreateComment(userID, &newComment)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, responses.CreateComment{
		ID:        comment.ID,
		Message:   comment.Message,
		PhotoID:   comment.PhotoID,
		UserID:    userID,
		CreatedAt: comment.CreatedAt,
	})
}

// GetComments godoc
// @Summary      Get comments
// @Description  Get comments.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Success      200  {object}  []responses.GetComment
// @Failure		 400 {object} responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /comments [get]
// @Security	 BearerAuth
func GetAllComments(ctx *gin.Context) {
	comments, err := database.GetAllComments()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	commentsResponse := make([]responses.GetComment, len(comments))
	users := make(map[uint]models.User)
	photos := make(map[uint]models.Photo)
	for i, comment := range comments {
		commentsResponse[i].Set(comment)
		user, ok := users[comment.UserID]
		if !ok {
			user, err = database.GetUserWithoutPreload(comment.UserID)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			users[comment.UserID] = user
		}
		commentsResponse[i].User = responses.UserComment{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}
		photo, ok := photos[comment.PhotoID]
		if !ok {
			photo, err = database.GetSinglePhoto(comment.PhotoID)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			photos[comment.PhotoID] = photo
		}
		commentsResponse[i].Photo = photo
	}
	ctx.JSON(http.StatusOK, commentsResponse)
}

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Update a comment associated with logged in user.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param		 commentId path uint true "ID number of the comment"
// @Param        comment body dto.CommentMessage true "New JSON of the comment."
// @Success      200  {object}  models.Comment
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /comments/{commentId} [put]
// @Security	 BearerAuth
func UpdateComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	parsedID, err := strconv.ParseUint(commentID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	var commentDto dto.CommentMessage
	if err := ctx.ShouldBindJSON(&commentDto); err != nil {
		abortBadRequest(err, ctx)
		return
	}
	if err := validate.Struct(&commentDto); err != nil {
		validationAbort(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	comment, err := database.UpdateComment(uint(parsedID), userID, &commentDto)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Comment with ID %d is not found.", parsedID),
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
	ctx.JSON(http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Delete a comment associated with logged in user.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param		 commentId path uint true "ID number of the comment to be deleted"
// @Success      200  {object}  responses.Message
// @Failure      400  {object}  responses.ErrorMessage
// @Failure      403  {object}  responses.ErrorMessage
// @Failure      404  {object}  responses.ErrorMessage
// @Failure      500  {object}  nil
// @Router       /comments/{commentId} [delete]
// @Security	 BearerAuth
func DeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	parsedID, err := strconv.ParseUint(commentID, 10, 0)
	if err != nil {
		abortBadRequest(err, ctx)
		return
	}
	userID, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := database.DeleteComment(uint(parsedID), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, responses.ErrorMessage{
				ErrorMessage: fmt.Sprintf("Comment with ID %d is not found.", parsedID),
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
