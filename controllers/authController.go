package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"example.com/m/social_media/db"
	"example.com/m/social_media/models"
	"example.com/m/social_media/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
//
//	@Summary		Creates a user and persist to database
//	@Description	Simple discription to function
//	@Tags			user
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{object}	models.User{}
//	@Router			/register [post]
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
	task, err := utils.NewEmailDeliveryTask(user.ID, "some:template:id")
	if err != nil {
		log.Printf("Could not create email task: %v", err)
	}
	go db.REDISCLIENT.Enqueue(task)
	db.DB.Create(&user)
	ctx.JSON(http.StatusCreated, gin.H{"userID:": user.ID})
}

// Signin		godoc
//
//	@Summary		Authenticates use and provides JWT
//	@Description	Authenticates user
//	@Param			email		path	string	true	"email for signin"
//	@Param			password	path	string	true	"password required"
//	@Tags			user
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{string}	string
//	@Router			/login [post]
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
		"email": userFound.Email,
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token"})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get the currently authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /user [get]
func GetUserProfile(ctx *gin.Context) {
	user, _ := ctx.Get("currentUser")

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}
