package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/entity"
	"net/http"
	"strconv"
)

func (h *Handler) createList(c *gin.Context) {
	userId := c.GetInt(userCtx)
	if userId == 0 {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input entity.TodoList
	if err := c.BindJSON(input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listId, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{"id": listId})
}

type getAllListsResp struct {
	Data []entity.TodoList
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId := c.GetInt(userCtx)
	if userId == 0 {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResp{
		Data: lists,
	})
}

func (h *Handler) getList(c *gin.Context) {
	userId := c.GetInt(userCtx)
	if userId == 0 {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetOne(userId, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
