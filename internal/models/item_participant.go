package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemParticipant struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ItemID     uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SplitRatio float64   `gorm:"type:decimal(5,2);default:1.0" json:"split_ratio"`
	CreatedAt  time.Time `json:"created_at"`

	// Relations
	Item BillItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ip *ItemParticipant) BeforeCreate(tx *gorm.DB) error {
	if ip.ID == uuid.Nil {
		ip.ID = uuid.New()
	}
	return nil
}