package controllers

import (
	"dbo-test/internal/dal"
	"dbo-test/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Summary		Create a new user
//	@Description	Create a new user with the provided email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			input	body		createUserReq	true	"User details"
//	@Success		200		{object}	successResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/user [post]
func CreateUser(c *gin.Context) {
	var input createUserReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	if err := dal.User.Create(&model.User{
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: fmt.Sprintf("cannot create user: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: successStatus,
	})
}
