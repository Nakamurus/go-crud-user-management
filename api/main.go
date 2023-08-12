package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"

	"github.com/nakamurus/go-user-management/handlers"
	"github.com/nakamurus/go-user-management/middleware"
	"github.com/nakamurus/go-user-management/util"
)

type App struct {
	DB     *gorm.DB
	JWTKey []byte
	RDB    *redis.Client
}

func main() {
	var app App

	app.JWTKey = []byte(os.Getenv("JWT_KEY"))

	app.DB = util.DBConnect()
	app.RDB = util.RedisClient()

	r := gin.Default()
	r.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SECRET"),
		ErrorFunc: func(c *gin.Context) {
			c.String(http.StatusBadRequest, "CSRF token mismatch")
			c.Abort()
		},
	}))
	rl := util.NewRateLimiter(5)
	r.Use(func(c *gin.Context) {
		rl.MiddleWare()
		c.Next()
	})

	uh := handlers.UserHandler(app.DB, app.JWTKey)
	ah := handlers.AuthHandlerInit(app.DB, app.JWTKey, app.RDB)
	m := middleware.NewMiddleware(app.JWTKey)
	authorized := r.Group("/me")
	authorized.Use(m.AuthenticateMiddleware())
	{
		authorized.PUT("/:id", uh.UpdateUserHandler())
		authorized.PUT("/:id/password", ah.ChangePasswordHandler())
		authorized.DELETE("/:id", uh.DeleteUserHandler())
		authorized.POST("/refresh-token", ah.RefreshTokenHandler())
		authorized.GET("/:id", uh.GetUserHandler())
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.POST("`/login", ah.LoginHandler())
	r.GET("/users", uh.ListUsersHandler())
	r.GET("/user/:id", m.AuthenticateMiddleware(), uh.GetUserHandler())
	r.POST("/user", uh.CreateUserHandler())

	r.Run(":8080")
}
