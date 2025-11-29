package service

import (
	"errors"
	"patungan-server/internal/dto/request"
	"patungan-server/internal/dto/response"
	"patungan-server/internal/models"
	"patungan-server/internal/repository"

	"github.com/google/uuid"
)

type BillService interface {
	Create(req *request.CreateBillRequest, userID uuid.UUID) (*response.BillResponse, error)
	GetByID(id uuid.UUID) (*response.BillResponse, error)
	GetAll(userID uuid.UUID) ([]response.BillResponse, error)
	Update(id uuid.UUID, req *request.UpdateBillRequest) (*response.BillResponse, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
}

type billService struct {
	billRepo repository.BillRepository
}

func NewBillService(billRepo repository.BillRepository) BillService {
	return &billService{
		billRepo: billRepo,
	}
}

func (s *billService) Create(req *request.CreateBillRequest, userID uuid.UUID) (*response.BillResponse, error) {
	bill := &models.Bill{
		Title: req.Title,
		Description: req.Description,
		TotalAmount: req.TotalAmount,
		Currency: req.Currency,
		SplitMethod: req.SplitMethod,
		CreatedBy: userID,
		Status: "pending",
	}

	if err := s.billRepo.Create(bill); err != nil {
		return nil, err
	}

	// TODO : Add participants and in Phase 2
	
	createdBill, err := s.billRepo.FindByID(bill.ID)
	if err != nil {
		return nil, err
	}

	return s.toBillResponse(createdBill), nil
}

func (s *billService) GetByID(id uuid.UUID) (*response.BillResponse, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	return s.toBillResponse(bill), nil
}

func (s *billService) GetAll(userID uuid.UUID) ([]response.BillResponse, error) {
	bills, err := s.billRepo.FindAll(userID)
	if err != nil {
		return nil, err
	}

	var billResponses []response.BillResponse
	for _, bill := range bills {
		billResponses = append(billResponses, *s.toBillResponse(&bill))
	}

	return billResponses, nil
}

func (s *billService) Update(id uuid.UUID, req *request.UpdateBillRequest) (*response.BillResponse, error) {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bill not found")
	}

	if req.Title != "" {
		bill.Title = req.Title
	}
	if req.Description != "" {
		bill.Description = req.Description
	}
	if req.TotalAmount != 0 {
		bill.TotalAmount = req.TotalAmount
	}
	if req.Status != "" {
		bill.Status = req.Status
	}

	if err := s.billRepo.Update(bill); err != nil {
		return nil, err
	}

	return s.toBillResponse(bill), nil
}

func (s *billService) Delete(id uuid.UUID, userID uuid.UUID) error {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return errors.New("bill not found")
	}

	if bill.CreatedBy != userID {
		return errors.New("unauthorized to delete this bill")
	}

	if err := s.billRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *billService) toBillResponse(bill *models.Bill) *response.BillResponse {
	resp := &response.BillResponse{
		ID:          bill.ID,
		Title:       bill.Title,
		Description: bill.Description,
		TotalAmount: bill.TotalAmount,
		Currency:    bill.Currency,
		SplitMethod: bill.SplitMethod,
		Status:      bill.Status,
		CreatedBy:   bill.CreatedBy,
		CreatedAt:   bill.CreatedAt,
		UpdatedAt:   bill.UpdatedAt,
		Creator: response.UserResponse{
			ID:       bill.Creator.ID,
			Email:    bill.Creator.Email,
			FullName: bill.Creator.FullName,
		},
	}

	for _, p := range bill.Participants {
		resp.Participants = append(resp.Participants, response.ParticipantResponse{
			ID:         p.ID,
			UserID:     p.UserID,
			AmountOwed: p.AmountOwed,
			AmountPaid: p.AmountPaid,
			IsSettled:  p.IsSettled,
			User: response.UserResponse{
				ID:       p.User.ID,
				Email:    p.User.Email,
				FullName: p.User.FullName,
			},
		})
	}

	return resp
}