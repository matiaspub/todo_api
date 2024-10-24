package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/service"
	"net/http"
	"strconv"
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

	// middlewares might be separate struct
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.GET("/", h.getAllLists)
			lists.POST("/", h.createList)
			lists.GET("/:id", h.getList)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.GET("/", h.getAllItems)
				items.POST("/", h.createItem)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItem)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}

func getParamInt(c *gin.Context, paramName string) (int, bool) {
	listId, err := strconv.Atoi(c.Param(paramName))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid "+paramName+" param")
		return 0, true
	}
	return listId, false
}

func currentUserId(c *gin.Context) (int, bool) {
	userId := c.GetInt(userCtx)
	if userId == 0 {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, true
	}
	return userId, false
}
