package server

import (
	"dbo-test/internal/controllers"
	"dbo-test/internal/middlewares"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	//auth routes
	r.POST("/login", controllers.LoginHandler)

	r.Use(middlewares.JWTAuthMiddleware(os.Getenv("JWT_SECRET")))

	//customer routes
	customerGroup := r.Group("/customer")
	customerGroup.POST("/", controllers.CreateCustomer)
	customerGroup.GET("/", controllers.GetMultipleCustomer)
	customerGroup.GET("/:id", controllers.GetSingleCustomer)
	customerGroup.PUT("/:id", controllers.UpdateCustomer)
	customerGroup.DELETE("/:id", controllers.DeleteCustomer)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
