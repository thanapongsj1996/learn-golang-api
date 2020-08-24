package migrations

import (
	"learn-golang-api/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1598041094CreateArticlesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1598041094",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("articles").Error
		},
	}
}
