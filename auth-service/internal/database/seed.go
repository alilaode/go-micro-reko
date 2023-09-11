package database

import (
	"auth-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	users := []model.User{
		{
			ID:       uuid.NewString(),
			Username: "alilaode",
			Hash:     "123",
		},
	}

	db.AutoMigrate(model.User{})

	if err := db.First(&model.User{}).Error; err == gorm.ErrRecordNotFound {
		db.Create(&users)
	}

}
