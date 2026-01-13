package repository

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/user/go-todo-api/internal/database"
	"github.com/user/go-todo-api/internal/models"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

// GenerateRefreshToken creates a cryptographically secure random token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Create stores a new refresh token
func (r *TokenRepository) Create(userID uint, token string, expiresAt time.Time) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	err := database.DB.Create(refreshToken).Error
	return refreshToken, err
}

// FindByToken retrieves a refresh token
func (r *TokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := database.DB.Where("token = ? AND revoked = ?", token, false).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// Revoke marks a refresh token as revoked
func (r *TokenRepository) Revoke(token string) error {
	return database.DB.Model(&models.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *TokenRepository) RevokeAllForUser(userID uint) error {
	return database.DB.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

// CleanupExpired removes expired tokens
func (r *TokenRepository) CleanupExpired() error {
	return database.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}
