package api

import (
	"bot/api/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGin(h *handler.HTTPHandler) *gin.Engine {
	r := gin.Default()

	// CORS konfiguratsiyasi
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	password := r.Group("/password")
	{
		password.POST("", h.CreatePassword)
		password.GET("/:phone", h.GetAllPasswordsByUserID)
		password.GET("", h.GetByName)
	}
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, url))

	return r
}
