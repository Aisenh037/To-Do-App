package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

func ValidationErrorResponse(c *gin.Context, message string) {
	c.JSON(400, APIResponse{
		Success: false,
		Error:   message,
	})
}

// AppError represents a custom application error
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// HandleError handles different types of errors and returns appropriate responses
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		ErrorResponse(c, appErr.Code, appErr.Message)
		return
	}

	// Default internal server error
	ErrorResponse(c, 500, "An unexpected error occurred")
}
