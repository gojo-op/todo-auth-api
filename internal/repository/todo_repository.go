package repository

import (
	"errors"

	"github.com/gojo-op/todo-auth-api/internal/models"
	"gorm.io/gorm"
)

var ErrTodoNotFound = errors.New("todo not found")

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) FindAllByUserID(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) FindByIDAndUserID(id, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *TodoRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTodoNotFound
	}
	return nil
}
