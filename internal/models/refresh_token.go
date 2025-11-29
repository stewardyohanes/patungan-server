package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Token     string         `gorm:"type:varchar(500);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	RevokedAt *time.Time     `gorm:"index" json:"revoked_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook untuk generate UUID
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

// IsValid checks if token is still valid
func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}

// Revoke marks the token as revoked
func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = &now
}
