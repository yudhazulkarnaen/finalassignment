package controllers

import (
	"net/http"

	"finalassignment.id/finalassignment/controllers/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const uniqueViolationErr = "23505"

func validationAbort(err error, ctx *gin.Context) {
	var errorMessage string
	for _, err := range err.(validator.ValidationErrors) {
		errorMessage += err.Error() + "\n"
	}
	if len(errorMessage) > 0 {
		errorMessage = errorMessage[:len(errorMessage)-1]
	}
	abortBadRequest(err, ctx)
}

func abortBadRequest(err error, ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.ErrorMessage{
		ErrorMessage: err.Error(),
	})
}
