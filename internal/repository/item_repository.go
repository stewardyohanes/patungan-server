package repository

import (
	"patungan-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *models.BillItem) error
	CreateBatch(items []models.BillItem) error
	FindByID(id uuid.UUID) (*models.BillItem, error)
	FindByBillID(billID uuid.UUID) ([]models.BillItem, error)
	Update(item *models.BillItem) error
	Delete(id uuid.UUID) error
	DeleteByBillID(billID uuid.UUID) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *models.BillItem) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) CreateBatch(items []models.BillItem) error {
	return r.db.Create(&items).Error
}

func (r *itemRepository) FindByID(id uuid.UUID) (*models.BillItem, error) {
	var item models.BillItem
	err := r.db.Preload("Participants.User").First(&item, "id = ?", id).Error
	return &item, err
}

func (r *itemRepository) FindByBillID(billID uuid.UUID) ([]models.BillItem, error) {
	var items []models.BillItem
	err := r.db.Preload("Participants.User").Where("bill_id = ?", billID).Find(&items).Error
	return items, err
}

func (r *itemRepository) Update(item *models.BillItem) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BillItem{}, "id = ?", id).Error
}

func (r *itemRepository) DeleteByBillID(billID uuid.UUID) error {
	return r.db.Delete(&models.BillItem{}, "bill_id = ?", billID).Error
}