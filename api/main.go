package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nakamurus/go-user-management/handlers"
	"github.com/nakamurus/go-user-management/middleware"
	"github.com/nakamurus/go-user-management/models"
	"github.com/nakamurus/go-user-management/util"
)

type App struct {
	DB     *gorm.DB
	JWTKey []byte
}

func main() {
	var app App

	app.JWTKey = []byte(os.Getenv("JWT_KEY"))

	connStr := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to retrieve the database connection")
	}
	defer sqlDB.Close()

	db.AutoMigrate(&models.User{})
	// create user table if not exists
	db.Exec(`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	DROP TABLE IF EXISTS users;

	CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4() NOT NULL,
		name varchar(255) NOT NULL,
		email varchar(255) NOT NULL,
		password varchar(255) NOT NULL,
		PRIMARY KEY (id)
	);
	`)

	app.DB = db

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
	ah := handlers.AuthHandlerInit(app.DB, app.JWTKey)
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
