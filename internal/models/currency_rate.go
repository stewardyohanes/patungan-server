package models

import (
	"time"

	"gorm.io/gorm"
)

type CurrencyRate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FromCurrency string    `gorm:"type:varchar(3);not null;index" json:"from_currency"`
	ToCurrency   string    `gorm:"type:varchar(3);not null;index" json:"to_currency"`
	Rate         float64   `gorm:"type:decimal(18,6);not null" json:"rate"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (cr *CurrencyRate) BeforeCreate(tx *gorm.DB) error {
	return nil
}