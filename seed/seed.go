package seed

import (
	"learn-golang-api/config"
	"learn-golang-api/migrations"
	"learn-golang-api/models"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/prometheus/common/log"
)

func Load() {
	db := config.GetDB()

	// Clean Database
	db.DropTableIfExists("articles", "categories", "migrations")
	migrations.Migrate()

	// Add categories
	log.Info("Creating categories...")

	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)

	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}

		db.Create(&category)
		categories = append(categories, category)
	}

	// Add articles
	log.Info("Creating articles...")

	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)

	for i := 1; i <= numOfArticles; i++ {
		article := models.Article{
			Title:      faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			CategoryID: uint(rand.Intn(numOfCategories) + 1),
		}

		db.Create(&article)
		articles = append(articles, article)
	}
}
