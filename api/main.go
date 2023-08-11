package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nakamurus/go-user-management/handlers"
	"github.com/nakamurus/go-user-management/models"
)

var DB *gorm.DB

func init() {
	var err error
	config := models.DBConfig{
		HOST:     os.Getenv("DB_HOST"),
		PORT:     os.Getenv("DB_PORT"),
		USER:     os.Getenv("DB_USER"),
		DBNAME:   os.Getenv("DB_NAME"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
	}
	connStr := "postgres://" + config.USER + ":" + config.PASSWORD + "@" +
		config.HOST + ":" + config.PORT + "/" + config.DBNAME + "?sslmode=disable"
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB.AutoMigrate(&models.User{})
	// create user table if not exists
	DB.Exec(`
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
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/users", handlers.ListUsersHandler(DB))
	r.GET("/user/:id", handlers.GetUserHandler(DB))
	r.POST("/user", handlers.CreateUserHandler(DB))
	r.PUT("/user/:id", handlers.UpdateUserHandler(DB))
	r.DELETE("/user/:id", handlers.DeleteUserHandler(DB))

	r.Run(":8080")
}
