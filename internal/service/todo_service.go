package service

import (
	"github.com/gojo-op/todo-auth-api/internal/models"
	"github.com/gojo-op/todo-auth-api/internal/repository"
)

type TodoService struct {
	todos *repository.TodoRepository
}

type CreateTodoInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTodoInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func NewTodoService(todos *repository.TodoRepository) *TodoService {
	return &TodoService{todos: todos}
}

func (s *TodoService) List(userID uint) ([]models.Todo, error) {
	return s.todos.FindAllByUserID(userID)
}

func (s *TodoService) Create(userID uint, input CreateTodoInput) (*models.Todo, error) {
	todo := &models.Todo{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
	}

	if err := s.todos.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) Get(userID, todoID uint) (*models.Todo, error) {
	return s.todos.FindByIDAndUserID(todoID, userID)
}

func (s *TodoService) Update(userID, todoID uint, input UpdateTodoInput) (*models.Todo, error) {
	todo, err := s.todos.FindByIDAndUserID(todoID, userID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		todo.Title = *input.Title
	}
	if input.Description != nil {
		todo.Description = *input.Description
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	if err := s.todos.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) Delete(userID, todoID uint) error {
	return s.todos.Delete(todoID, userID)
}
