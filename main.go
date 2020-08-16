package main

import (
	"learn-golang-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.Serve(r)
	r.Run()
}
