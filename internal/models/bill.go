package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	TotalAmount float64 `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Currency string `gorm:"type:varchar(3);default:'IDR'" json:"currency"`
	SplitMethod string `gorm:"type:varchar(20);not null" json:"split_method"` // equal, items, custom
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Status string `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, settled, cancelled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Participants []BillParticipant `gorm:"foreignKey:BillID" json:"participants,omitempty"`
	Items []BillItem `gorm:"foreignKey:BillID" json:"items,omitempty"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}