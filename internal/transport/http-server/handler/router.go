package handler

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/TTekmii/todo-list-app/internal/transport/http-server/middleware"

	_ "github.com/TTekmii/todo-list-app/docs"
)

func (h *Handler) InitRoutes(logger *slog.Logger) *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(gin.Recovery())

	router.Use(middleware.LoggingMiddleware(logger))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")

	api.Use(middleware.AuthMiddleware(h.services.Auth))
	{
		lists := api.Group("/lists")
		{
			lists.POST("", h.createList)
			lists.GET("", h.getAllLists)

			listsByID := lists.Group("/:id")
			{
				listsByID.GET("", h.getListById)
				listsByID.PUT("", h.updateList)
				listsByID.DELETE("", h.deleteList)

				items := listsByID.Group("/items")
				{
					items.POST("", h.createItem)
					items.GET("", h.getAllItems)

					itemsByID := items.Group("/:id")
					{
						itemsByID.GET("", h.getItemById)
						itemsByID.PUT("", h.updateItem)
						itemsByID.DELETE("", h.deleteItem)
					}
				}
			}
		}
	}
	return router
}
