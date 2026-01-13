package repository

import (
	"github.com/user/go-todo-api/internal/database"
	"github.com/user/go-todo-api/internal/models"
)

type TodoRepository struct{}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{}
}

// QueryParams holds pagination, filtering, and sorting options
type QueryParams struct {
	Page     int
	PageSize int
	Status   string
	Search   string
	SortBy   string
	SortDir  string
}

// PaginatedResult holds the paginated response
type PaginatedResult struct {
	Data       []models.Todo `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return database.DB.Create(todo).Error
}

func (r *TodoRepository) FindAllByUserID(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&todos).Error
	return todos, err
}

// FindAllWithFilters returns paginated, filtered, and sorted todos
func (r *TodoRepository) FindAllWithFilters(userID uint, params QueryParams) (*PaginatedResult, error) {
	var todos []models.Todo
	var total int64

	// Base query
	query := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID)

	// Apply status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Apply search filter (search in title and description)
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := params.SortDir
	if sortDir == "" {
		sortDir = "DESC"
	}
	// Validate sort fields to prevent SQL injection
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"title":      true,
		"status":     true,
		"due_date":   true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}
	query = query.Order(sortBy + " " + sortDir)

	// Apply pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResult{
		Data:       todos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *TodoRepository) FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := database.DB.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) FindByIDAndUserID(id, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return database.DB.Save(todo).Error
}

func (r *TodoRepository) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{}).Error
}
