package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/m/social_media/db"
	"example.com/m/social_media/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *gin.Context) {
	var authInput models.AuthInput
	
	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse auth input"})
		return
	}
	
	var userFound models.User
	db.DB.Where("email=?", authInput.Email).Find(&userFound)
	if userFound.ID != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	db.DB.Create(&user)
	ctx.JSON(http.StatusCreated, gin.H{"user:": user})
}

func Signin(ctx *gin.Context) {
	var authInput models.AuthInput
	var userFound models.User

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Where("email=?", authInput.Email).Find(&userFound)
	if userFound.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userFound.ID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token"})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func GetUserProfile(ctx *gin.Context) {
	user, _ := ctx.Get("currentUser")

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}
