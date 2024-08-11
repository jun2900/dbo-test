package main

import (
	"dbo-test/internal/server"
	"fmt"

	_ "dbo-test/docs"
)

// @title						DBO-TEST API
// @version					1.0
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
