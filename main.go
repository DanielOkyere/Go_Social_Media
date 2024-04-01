package main

import (
	"net/http"

	"example.com/m/social_media/db"
	_ "example.com/m/social_media/docs"
	"example.com/m/social_media/routes"
	"example.com/m/social_media/utils"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

//	@title			Go Social Media API
//	@version		1.0
//	@description	This project is social media api

//	@contact.name	Daniel Okyere
//	@contact.url	http://www.swagger.io/support
//	@contact.email	daniel.kwame.okyere101@gmail.com

//	@host	localhost:3000

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https:://swagger.io/resources/open-api/

func main() {
	server := gin.Default()

	utils.LoadEnvs()
	db.InitDB()
	db.InitREDIS()
	go db.HandleServer()
	routes.RegisterRoutes(server)

	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I am alive"})
	})

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run()
}
