package routes

import (
	"example.com/m/social_media/controllers"
	"example.com/m/social_media/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.Use(middlewares.JSONMiddleware())
	server.POST("/register", controllers.CreateUser)
    server.POST("/login", controllers.Signin)

	v1:=server.Group("/user")
	{
		v1.Use(middlewares.CheckAuth)
		v1.GET("/me", controllers.GetUserProfile)
	}

	
}