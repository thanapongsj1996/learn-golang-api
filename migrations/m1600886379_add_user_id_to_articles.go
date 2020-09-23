package migrations

import (
	"learn-golang-api/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1600886379AddUserIDToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1600886379",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("user_id").Error
		},
	}
}
