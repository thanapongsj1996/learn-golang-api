package migrations

import (
	"learn-golang-api/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1598041094CreateArticlesTable(),
			m1600445650CreateCategoriesTable(),
			m1600620099AddCategoryIDToArticles(),
			m1600790848CreateUsersTable(),
			m1600886379AddUserIDToArticles(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}
