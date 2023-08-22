package database

import "github.com/Prameesh-P/user-handler/internal/user"

func SycnDB() {
	DB.AutoMigrate(user.User{})
}