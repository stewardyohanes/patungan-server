package handler

import (
	"patungan-server/internal/dto/request"
	"patungan-server/internal/service"
	"patungan-server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BillHandler struct {
	billService service.BillService
	validator   *validator.Validate
}

func NewBillHandler(billService service.BillService) *BillHandler {
	return &BillHandler{
		billService: billService,
		validator: validator.New(),
	}
}

func (h *BillHandler) Create(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req request.CreateBillRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

		if err := h.validator.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation failed", err.Error())
	}

	result, err := h.billService.Create(&req, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "bill created successfully", result)
}

func (h *BillHandler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	result, err := h.billService.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "bill retrieved successfully", result)
}

func (h *BillHandler) GetAll(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	result, err := h.billService.GetAll(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "bills retrieved successfully", result)
}

func (h *BillHandler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	var req request.UpdateBillRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	result, err := h.billService.Update(id, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "bill updated successfully", result)
}

func (h *BillHandler) Delete(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	if err := h.billService.Delete(id, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "bill deleted successfully", nil)
}