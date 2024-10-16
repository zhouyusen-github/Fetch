package router

import (
	"Fetch/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/add", handlers.AddPoints)
	router.POST("/spend", handlers.SpendPoints)
	router.GET("/balance", handlers.GetPointsBalance)

	return router
}
