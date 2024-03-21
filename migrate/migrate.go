package main

import (
	"example.com/m/social_media/db"
	"example.com/m/social_media/models"
	"example.com/m/social_media/utils"
)

func init() {
	utils.LoadEnvs()
	db.InitDB()
}

func main() {
	db.DB.AutoMigrate(&models.User{})
}
