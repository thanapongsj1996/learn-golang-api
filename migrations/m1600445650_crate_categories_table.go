package migrations

import (
	"learn-golang-api/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1600445650CreateCategoriesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1600445650",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("categories").Error
		},
	}
}
