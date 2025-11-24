package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillParticipant struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BillID uuid.UUID `gorm:"type:uuid;not null" json:"bill_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AmountOwed float64 `gorm:"type:decimal(12,2);not null" json:"amount_owed"`
	AmountPaid float64 `gorm:"type:decimal(12,2);default:0" json:"amount_paid"`
	IsSettled bool `gorm:"default:false" json:"is_settled"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Bill Bill `gorm:"foreignKey:BillID" json:"bill,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (bp *BillParticipant) BeforeCreate(tx *gorm.DB) error {
	if bp.ID == uuid.Nil {
		bp.ID = uuid.New()
	}
	return nil
}