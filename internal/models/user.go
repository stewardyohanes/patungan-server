package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	FullName string `gorm:"type:varchar(255);not null" json:"full_name"`
	Phone string `gorm:"type:varchar(50)" json:"phone"`
	AvatarURL string `gorm:"type:text" json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}