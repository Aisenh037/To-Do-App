package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken stores refresh tokens in the database
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"-" gorm:"unique;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time      `json:"expires_at"`
	Revoked   bool           `json:"revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds until access token expires
}

// RefreshTokenRequest for token refresh endpoint
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
