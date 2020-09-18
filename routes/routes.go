package routes

import (
	"learn-golang-api/config"
	"learn-golang-api/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	articlesGroup := r.Group("/api/v1/articles")
	articleController := controllers.Articles{DB: db}

	{
		articlesGroup.GET("/", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("/", articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}
}
