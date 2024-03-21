package routes

import (
	"example.com/m/social_media/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", controllers.CreateUser)
    server.POST("/login", controllers.Signin)
}