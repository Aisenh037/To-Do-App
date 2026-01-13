package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/go-todo-api/internal/config"
	"github.com/user/go-todo-api/internal/models"
	"github.com/user/go-todo-api/internal/repository"
	"github.com/user/go-todo-api/internal/worker"
	"github.com/user/go-todo-api/pkg/utils"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo:  repository.NewUserRepository(),
		tokenRepo: repository.NewTokenRepository(),
	}
}

// createTokenPair generates both access and refresh tokens
func (h *AuthHandler) createTokenPair(user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessToken, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := repository.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token (7 days expiry)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = h.tokenRepo.Create(user.ID, refreshTokenStr, expiresAt)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    config.AppConfig.JWTExpiryHours * 3600,
	}, nil
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Create a new user account and return access/refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RegisterRequest  true  "Registration Info"
// @Success      201      {object}  utils.APIResponse{data=map[string]interface{}}
// @Failure      400      {object}  utils.APIResponse
// @Failure      409      {object}  utils.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input: "+err.Error())
		return
	}

	// Check if email already exists
	if h.userRepo.ExistsByEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token pair
	tokenPair, err := h.createTokenPair(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Enqueue welcome email in background
	worker.GlobalWorker.Enqueue(worker.Task{
		Type: "SEND_WELCOME_EMAIL",
		Payload: map[string]interface{}{
			"email": user.Email,
			"name":  user.Name,
		},
	})

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"user":   user.ToResponse(),
		"tokens": tokenPair,
	})
}

// Login handles user authentication
// @Summary      Login user
// @Description  Authenticate user and return access/refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login Credentials"
// @Success      200      {object}  utils.APIResponse{data=map[string]interface{}}
// @Failure      400      {object}  utils.APIResponse
// @Failure      401      {object}  utils.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input: "+err.Error())
		return
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token pair
	tokenPair, err := h.createTokenPair(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user":   user.ToResponse(),
		"tokens": tokenPair,
	})
}

// RefreshToken issues a new access token using a refresh token
// @Summary      Refresh token
// @Description  Issue a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.RefreshTokenRequest  true  "Refresh Token"
// @Success      200      {object}  utils.APIResponse{data=models.TokenPair}
// @Failure      401      {object}  utils.APIResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid input: "+err.Error())
		return
	}

	// Find refresh token
	refreshToken, err := h.tokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Check if expired
	if time.Now().After(refreshToken.ExpiresAt) {
		h.tokenRepo.Revoke(req.RefreshToken)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// Get user
	user, err := h.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	// Revoke old refresh token
	h.tokenRepo.Revoke(req.RefreshToken)

	// Generate new token pair
	tokenPair, err := h.createTokenPair(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed", tokenPair)
}

// Logout revokes all refresh tokens for the user
// @Summary      Logout user
// @Description  Revoke all refresh tokens for the authenticated user
// @Tags         auth
// @Security     Bearer
// @Produce      json
// @Success      200      {object}  utils.APIResponse
// @Failure      401      {object}  utils.APIResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Revoke all refresh tokens for user
	h.tokenRepo.RevokeAllForUser(userID.(uint))

	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

// GetProfile returns the authenticated user's profile
// @Summary      Get user profile
// @Description  Get the profile information of the authenticated user
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200      {object}  utils.APIResponse{data=models.UserResponse}
// @Failure      401      {object}  utils.APIResponse
// @Failure      404      {object}  utils.APIResponse
// @Router       /profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved", user.ToResponse())
}
