package repository

import (
	"patungan-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillRepository interface {
	Create(bill *models.Bill) error
	FindByID(id uuid.UUID) (*models.Bill, error)
	FindAll(userID uuid.UUID) ([]models.Bill, error)
	Update(bill *models.Bill) error
	Delete(id uuid.UUID) error
}

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) Create(bill *models.Bill) error {
	return r.db.Create(bill).Error
}

func (r *billRepository) FindByID(id uuid.UUID) (*models.Bill, error) {
	var bill models.Bill
	err := r.db.Preload("Creator").Preload("Participants").Preload("Items").First(&bill, "id = ?", id).Error
	return &bill, err
}

func (r *billRepository) FindAll(userID uuid.UUID) ([]models.Bill, error) {
	var bills []models.Bill
	err := r.db.
		Preload("Creator").
		Preload("Participants").
		Joins("LEFT JOIN bill_participants ON bills.id = bill_participants.bill_id").
		Where("bills.created_by = ? OR bill_participants.user_id = ?", userID, userID).
		Group("bills.id").
		Find(&bills).Error
	return bills, err
}

func (r *billRepository) Update(bill *models.Bill) error {
	return r.db.Save(bill).Error
}

func (r *billRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Bill{}, "id = ?", id).Error
}