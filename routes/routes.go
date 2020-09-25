package routes

import (
	"learn-golang-api/config"
	"learn-golang-api/controllers"
	"learn-golang-api/middlewares"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")

	// เป็นการอ่านค่าของ JWT จาก header แล้วเอา JWT มาดึงในส่วน payload
	// หาค่าของ sub แล้วเอา sub ไปหา user แล้วคืน user กลับมาในส่วน context
	authenticate := middlewares.Authenticate().MiddlewareFunc()

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middlewares.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	// usersGroup := v1.Group("users")
	// usersController := controllers.Users{DB: db}
	// {
	// 	usersGroup.GET("", usersController.FindAll)
	// 	usersGroup.POST("", usersController.Create)
	// 	usersGroup.GET("/:id", usersController.FindOne)
	// 	usersGroup.PATCH("/:id", usersController.Update)
	// 	usersGroup.DELETE("/:id", usersController.Delete)
	// 	usersGroup.PATCH("/:id/promote", usersController.Promote)
	// 	usersGroup.PATCH("/:id/demote", usersController.Demote)
	// }

	articlesGroup := v1.Group("articles")
	articleController := controllers.Articles{DB: db}
	{
		articlesGroup.GET("/", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("/", authenticate, articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}

	categoriesGroup := v1.Group("categories")
	categoryController := controllers.Categories{DB: db}
	{
		categoriesGroup.GET("/", categoryController.FindAll)
		categoriesGroup.GET("/:id", categoryController.FindOne)
		categoriesGroup.POST("/", categoryController.Create)
		categoriesGroup.PATCH("/:id", categoryController.Update)
		categoriesGroup.DELETE("/:id", categoryController.Delete)
	}
}
