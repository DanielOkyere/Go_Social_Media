package main

import (
	"net/http"

	"example.com/m/social_media/db"
	"example.com/m/social_media/routes"
	"example.com/m/social_media/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	utils.LoadEnvs()
	db.InitDB()
	routes.RegisterRoutes(server)

	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I am alive"})
	})

	server.Run()
}
