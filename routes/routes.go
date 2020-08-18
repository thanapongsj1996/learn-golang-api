package routes

import (
	"learn-golang-api/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {

	articlesGroup := r.Group("/api/v1/articles")
	articleController := controllers.Articles{}

	{
		articlesGroup.GET("/", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("/", articleController.Create)
	}
}
