package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/entity"
	"net/http"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, done := currentUserId(c)
	if !done {
		return
	}

	listId, done := getParamInt(c, "id")
	if !done {
		return
	}

	var item entity.TodoItem
	if err := c.BindJSON(item); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, item)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]int{"id": id})
}

type getAllItemsResp struct {
	Data []entity.TodoItem
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, done := currentUserId(c)
	if !done {
		return
	}

	listId, done := getParamInt(c, "id")
	if !done {
		return
	}

	list, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllItemsResp{
		Data: list,
	})
}

func (h *Handler) getItem(c *gin.Context) {
	userId, done := currentUserId(c)
	if !done {
		return
	}

	itemId, done := getParamInt(c, "id")
	if !done {
		return
	}

	item, err := h.services.TodoItem.GetOne(userId, itemId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, done := currentUserId(c)
	if !done {
		return
	}

	itemId, done := getParamInt(c, "id")
	if !done {
		return
	}

	var input entity.UpdateTodoItemInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoItem.Update(userId, itemId, input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, done := currentUserId(c)
	if !done {
		return
	}

	itemId, done := getParamInt(c, "id")
	if !done {
		return
	}

	err := h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
