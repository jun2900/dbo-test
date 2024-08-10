package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	successStatus = "success"
	errorStatus   = "error"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type successResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type orderBy struct {
	Field string
	Desc  bool
}

type PagedResults struct {
	Page         int64       `json:"page"`
	PageSize     int64       `json:"page_size"`
	Data         interface{} `json:"data"`
	TotalRecords int         `json:"total_records"`
}

func generateJWT(email string) (string, error) {
	expire, err := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(expire)).Unix()
	claims["email"] = email

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
