package migrations

import (
	"learn-golang-api/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1600790848CreateUsersTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1600790848",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("users").Error
		},
	}
}
