package controllers

import (
	"dbo-test/internal/dal"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary		Logs in a user
// @Description	Logs in a user with email and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			input	body		loginReq	true	"login req"
// @Success		200		{object}	successResponse
// @Failure		400		{object}	errorResponse
// @Router			/login [post]
func LoginHandler(c *gin.Context) {
	var input loginReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	user, err := dal.User.Where(dal.User.Email.Eq(input.Email)).First()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid user",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid user",
		})
		return
	}

	accessToken, err := generateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: fmt.Sprintf("cannot generate jwt: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: successStatus,
		Data:   gin.H{"access_token": accessToken},
	})
}
