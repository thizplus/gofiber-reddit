package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications retrieves all notifications for the user
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	notifications, err := h.notificationService.GetNotifications(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve notifications", err)
	}

	return utils.SuccessResponse(c, "Notifications retrieved successfully", notifications)
}

// GetUnreadNotifications retrieves unread notifications
func (h *NotificationHandler) GetUnreadNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	notifications, err := h.notificationService.GetUnreadNotifications(c.Context(), userID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve unread notifications", err)
	}

	return utils.SuccessResponse(c, "Unread notifications retrieved successfully", notifications)
}

// GetNotification retrieves a single notification
func (h *NotificationHandler) GetNotification(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid notification ID")
	}

	notification, err := h.notificationService.GetNotification(c.Context(), notificationID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Notification not found", err)
	}

	return utils.SuccessResponse(c, "Notification retrieved successfully", notification)
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid notification ID")
	}

	err = h.notificationService.MarkAsRead(c.Context(), notificationID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark notification as read", err)
	}

	return utils.SuccessResponse(c, "Notification marked as read", nil)
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.notificationService.MarkAllAsRead(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark all notifications as read", err)
	}

	return utils.SuccessResponse(c, "All notifications marked as read", nil)
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid notification ID")
	}

	err = h.notificationService.DeleteNotification(c.Context(), notificationID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete notification", err)
	}

	return utils.SuccessResponse(c, "Notification deleted successfully", nil)
}

// DeleteAllNotifications deletes all notifications
func (h *NotificationHandler) DeleteAllNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.notificationService.DeleteAllNotifications(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete all notifications", err)
	}

	return utils.SuccessResponse(c, "All notifications deleted successfully", nil)
}

// GetUnreadCount retrieves unread notification count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	count, err := h.notificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve unread count", err)
	}

	return utils.SuccessResponse(c, "Unread count retrieved successfully", fiber.Map{
		"count": count,
	})
}

// GetSettings retrieves notification settings
func (h *NotificationHandler) GetSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	settings, err := h.notificationService.GetSettings(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve notification settings", err)
	}

	return utils.SuccessResponse(c, "Notification settings retrieved successfully", settings)
}

// UpdateSettings updates notification settings
func (h *NotificationHandler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.NotificationSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	settings, err := h.notificationService.UpdateSettings(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update notification settings", err)
	}

	return utils.SuccessResponse(c, "Notification settings updated successfully", settings)
}
