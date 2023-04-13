package initializers

import "github.com/nickemma/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}