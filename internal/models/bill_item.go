package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BillID    uuid.UUID `gorm:"type:uuid;not null" json:"bill_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Price     float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Quantity  int       `gorm:"default:1" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Bill Bill `gorm:"foreignKey:BillID" json:"bill,omitempty"`
	Participants []ItemParticipant `gorm:"foreignKey:ItemID" json:"participants,omitempty"`
}

func (bi *BillItem) BeforeCreate(tx *gorm.DB) error {
	if bi.ID == uuid.Nil {
		bi.ID = uuid.New()
	}
	return nil
}