package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nakamurus/go-user-management/handlers"
	"github.com/nakamurus/go-user-management/middleware"
	"github.com/nakamurus/go-user-management/models"
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
	uh := handlers.UserHandler(app.DB, app.JWTKey)
	ah := handlers.AuthHandlerInit(app.DB, app.JWTKey)
	m := middleware.NewMiddleware(app.JWTKey)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.POST("`/login", ah.LoginHandler())
	r.GET("/protected", m.AuthenticateMiddleware(), uh.GetUserHandler())
	r.GET("/users", uh.ListUsersHandler())
	r.GET("/user/:id", uh.GetUserHandler())
	r.POST("/user", uh.CreateUserHandler())
	r.POST("/refresh-token", ah.RefreshTokenHandler())
	r.PUT("/user/:id", uh.UpdateUserHandler())
	r.PUT("/user/:id/password", ah.ChangePasswordHandler())
	r.DELETE("/user/:id", uh.DeleteUserHandler())

	r.Run(":8080")
}
