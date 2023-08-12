package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nakamurus/go-user-management/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	JWTKey []byte
}

func UserHandler(db *gorm.DB, jwtKey []byte) *Handler {
	return &Handler{
		db:     db,
		JWTKey: jwtKey,
	}
}

func getUUIDFromRequest(c *gin.Context) uuid.UUID {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		panic(err)
	}
	return uid
}

func (h *Handler) ListUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := models.GetAllUsers(h.db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := getUUIDFromRequest(c)
		user, err := models.GetUserById(h.db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) CreateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := models.CreateUser(h.db, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user created" + newUser.Name})
	}
}

func (h *Handler) UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := getUUIDFromRequest(c)
		var user models.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.ID = id
		updatedUser, err := models.UpdateUser(h.db, &user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *Handler) DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := getUUIDFromRequest(c)

		user, err := models.GetUserById(h.db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := models.DeleteUser(h.db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}
