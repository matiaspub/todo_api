package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		lists := api.Group("/lists")
		{
			lists.GET("/", h.readAllLists)
			lists.POST("/", h.createList)
			lists.GET("/:id", h.readList)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.GET("/", h.readAllItems)
				items.POST("/", h.createItem)
				items.GET("/:itemID", h.readItem)
				items.PUT("/:itemID", h.updateItem)
				items.DELETE("/:itemID", h.deleteItem)
			}
		}
	}

	return router
}
