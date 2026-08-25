package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gojo-op/todo-auth-api/internal/repository"
	"github.com/gojo-op/todo-auth-api/internal/service"
)

type TodoHandler struct {
	todos *service.TodoService
}

func NewTodoHandler(todos *service.TodoService) *TodoHandler {
	return &TodoHandler{todos: todos}
}

func (h *TodoHandler) List(c *gin.Context) {
	userID := c.GetUint("userID")

	todos, err := h.todos.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var input service.CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todos.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) Get(c *gin.Context) {
	userID := c.GetUint("userID")
	todoID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := h.todos.Get(userID, todoID)
	if errors.Is(err, repository.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	todoID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var input service.UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todos.Update(userID, todoID, input)
	if errors.Is(err, repository.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	todoID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	if err := h.todos.Delete(userID, todoID); errors.Is(err, repository.ErrTodoNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}
