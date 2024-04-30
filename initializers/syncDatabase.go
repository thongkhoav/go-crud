package initializers

import (
	"github.com/thongkhoav/go-crud/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.Post{})
	DB.AutoMigrate(&models.User{})
}
