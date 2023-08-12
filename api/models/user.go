package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nakamurus/go-user-management/util"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type DBConfig struct {
	HOST     string
	PORT     string
	USER     string
	DBNAME   string
	PASSWORD string
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := util.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

func CreateUser(db *gorm.DB, user User) (*User, error) {
	fmt.Println("Start Create User")
	fmt.Println(user)
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetUserById(db *gorm.DB, id uuid.UUID) (*User, error) {
	var user User
	result := db.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *User) (*User, error) {
	result := db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func DeleteUser(db *gorm.DB, user *User) error {
	result := db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
