package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/entity"
	"net/http"
)

// @Summary Create Todo List
// @Tags lists
// @Description create todo list
// @Accept json
// @Produce json
// @Param input body entity.TodoList true "list info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} Error
// @Failure 500 {object} Error
// @Failure default {object} Error
// @Security ApiKeyAuth
// @Router /api/list [post]
func (h *Handler) createList(c *gin.Context) {
	userId, done := currentUserId(c)
	if done {
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
	userId, done := currentUserId(c)
	if done {
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
	userId, done := currentUserId(c)
	if done {
		return
	}

	listId, done2 := getParamInt(c, "id")
	if done2 {
		return
	}

	list, err := h.services.TodoList.GetOne(userId, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, done := currentUserId(c)
	if done {
		return
	}

	listId, done2 := getParamInt(c, "id")
	if done2 {
		return
	}

	var todoList entity.UpdateListInput
	if err := c.BindJSON(todoList); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err := h.services.TodoList.Update(userId, listId, todoList)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, done := currentUserId(c)
	if done {
		return
	}

	listId, done2 := getParamInt(c, "id")
	if done2 {
		return
	}

	err := h.services.TodoList.Delete(userId, listId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
