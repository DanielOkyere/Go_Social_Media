package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/m/social_media/db"
	"example.com/m/social_media/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 || authToken[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	db.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access denied"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("currentUser", user)
	ctx.Next()
}
