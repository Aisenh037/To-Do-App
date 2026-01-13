package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/go-todo-api/internal/middleware"
	"github.com/user/go-todo-api/internal/models"
	"github.com/user/go-todo-api/internal/repository"
	"github.com/user/go-todo-api/internal/worker"
	"github.com/user/go-todo-api/pkg/utils"
)

type TodoHandler struct {
	todoRepo *repository.TodoRepository
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todoRepo: repository.NewTodoRepository(),
	}
}

// GetAll returns all todos for the authenticated user with pagination, filtering, and sorting
// @Summary      Get all todos
// @Description  Get a paginated list of todos for the authenticated user with optional filtering and sorting
// @Tags         todos
// @Security     Bearer
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        page_size  query     int     false  "Items per page (default: 10, max: 100)"
// @Param        status     query     string  false  "Filter by status (pending, in_progress, completed)"
// @Param        search     query     string  false  "Search in title and description"
// @Param        sort_by    query     string  false  "Sort field (created_at, title, status, etc.)"
// @Param        sort_dir   query     string  false  "Sort direction (ASC, DESC)"
// @Success      200        {object}  utils.APIResponse{data=map[string]interface{}}
// @Failure      401        {object}  utils.APIResponse
// @Router       /todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "DESC")

	params := repository.QueryParams{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}

	result, err := h.todoRepo.FindAllWithFilters(userID, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}

	// Convert to response DTOs
	var todosResponse []models.TodoResponse
	for _, todo := range result.Data {
		todosResponse = append(todosResponse, todo.ToResponse())
	}

	utils.SuccessResponse(c, http.StatusOK, "Todos retrieved", gin.H{
		"todos":       todosResponse,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetByID returns a single todo by ID
// @Summary      Get todo by ID
// @Description  Get detailed information about a specific todo item
// @Tags         todos
// @Security     Bearer
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  utils.APIResponse{data=models.TodoResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid todo ID")
		return
	}

	todo, err := h.todoRepo.FindByIDAndUserID(uint(todoID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Todo not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo retrieved", todo.ToResponse())
}

// Create creates a new todo
// @Summary      Create a todo
// @Description  Create a new todo item for the authenticated user
// @Tags         todos
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateTodoRequest  true  "Todo Information"
// @Success      201      {object}  utils.APIResponse{data=models.TodoResponse}
// @Failure      400      {object}  utils.APIResponse
// @Failure      401      {object}  utils.APIResponse
// @Router       /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)

	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input: "+err.Error())
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = models.StatusPending
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     req.DueDate,
		UserID:      userID,
	}

	if err := h.todoRepo.Create(todo); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Todo created", todo.ToResponse())
}

// Update updates an existing todo
// @Summary      Update a todo
// @Description  Update the details of an existing todo item
// @Tags         todos
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Todo ID"
// @Param        request  body      models.UpdateTodoRequest  true  "Updated Todo Info"
// @Success      200      {object}  utils.APIResponse{data=models.TodoResponse}
// @Failure      400      {object}  utils.APIResponse
// @Failure      401      {object}  utils.APIResponse
// @Failure      404      {object}  utils.APIResponse
// @Router       /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid todo ID")
		return
	}

	// Find existing todo
	todo, err := h.todoRepo.FindByIDAndUserID(uint(todoID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Todo not found")
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input: "+err.Error())
		return
	}

	// Update fields if provided
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Status != "" {
		todo.Status = req.Status
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := h.todoRepo.Update(todo); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	// Enqueue notification if todo is completed
	if todo.Status == "completed" {
		worker.GlobalWorker.Enqueue(worker.Task{
			Type: "TODO_COMPLETED_NOTIFICATION",
			Payload: map[string]interface{}{
				"id":    todo.ID,
				"title": todo.Title,
			},
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo updated", todo.ToResponse())
}

// Delete deletes a todo
// @Summary      Delete a todo
// @Description  Remove a specific todo item permanently
// @Tags         todos
// @Security     Bearer
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid todo ID")
		return
	}

	// Verify todo exists and belongs to user
	_, err = h.todoRepo.FindByIDAndUserID(uint(todoID), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Todo not found")
		return
	}

	if err := h.todoRepo.Delete(uint(todoID), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo deleted", nil)
}
