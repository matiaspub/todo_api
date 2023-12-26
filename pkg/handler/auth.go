package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matiaspub/todo-api/pkg/entity"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) signIn(c *gin.Context) {

}
