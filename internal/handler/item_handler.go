package handler

import (
	"patungan-server/internal/dto/request"
	"patungan-server/internal/service"
	"patungan-server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ItemHandler struct {
	itemService service.ItemService
	validator   *validator.Validate
}

func NewItemHandler(itemService service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		validator: validator.New(),
	}
}

func (h *ItemHandler) AddItems(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	var req request.AddItemsRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation failed", err.Error())
	}

	if err := h.itemService.AddItems(billID, &req, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "items added successfully", nil)
}

func (h *ItemHandler) GetItems(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	result, err := h.itemService.GetItems(billID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "items fetched successfully", result)
}

func (h *ItemHandler) UpdateItem(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	itemID, err := uuid.Parse(c.Params("itemId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid item id", nil)
	}

	var req request.UpdateItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	result, err := h.itemService.UpdateItem(itemID, &req, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "item updated sucessfully", result)
}

func (h *ItemHandler) DeleteItem(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	itemID, err := uuid.Parse(c.Params("itemId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid item id", nil)
	}

	if err := h.itemService.DeleteItem(itemID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "item deleted successfully", nil)
}

