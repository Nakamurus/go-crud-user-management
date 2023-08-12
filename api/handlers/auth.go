package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakamurus/go-user-management/models"
	"github.com/nakamurus/go-user-management/util"
	"gorm.io/gorm"
)

var jwtkey = []byte(os.Getenv("JWT_SECRET"))

func generateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.Email,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func refreshJWTToken(c *gin.Context, user models.User) {
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var foundUser models.User

		if err := db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// ユーザーが入力したパスワードと、データベースに保存されているハッシュ化されたパスワードを比較
		if !util.CheckPasswordHash(user.Password, foundUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		tokenString, err := generateToken(foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func ChangePassword(db *gorm.DB, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := models.GetUserById(db, userID)

	if err != nil {
		return err
	}

	if !util.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("current password is wrong")
	}

	user.Password = newPassword
	return db.Save(&user).Error

}

func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
