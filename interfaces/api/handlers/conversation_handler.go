package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	chatWebsocket "gofiber-template/infrastructure/websocket"
	"gofiber-template/pkg/utils"
)

type ConversationHandler struct {
	conversationService services.ConversationService
	conversationRepo    repositories.ConversationRepository
	chatHub             *chatWebsocket.ChatHub
}

func NewConversationHandler(conversationService services.ConversationService, conversationRepo repositories.ConversationRepository, chatHub *chatWebsocket.ChatHub) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		conversationRepo:    conversationRepo,
		chatHub:             chatHub,
	}
}

// GetOrCreateConversation creates or retrieves a 1-on-1 conversation
// GET /conversations/with/:username
func (h *ConversationHandler) GetOrCreateConversation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get username from URL path parameter
	username := c.Params("username")
	if username == "" {
		return utils.ValidationErrorResponse(c, "Username is required")
	}

	conversation, err := h.conversationService.GetOrCreateConversation(c.Context(), userID, username)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create conversation", err)
	}

	return utils.SuccessResponse(c, "Conversation retrieved successfully", conversation)
}

// ListConversations retrieves all conversations for the current user
// GET /conversations?cursor=xxx&limit=20
func (h *ConversationHandler) ListConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get cursor from query params
	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Get limit from query params (default: 20, max: 50)
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}
	}

	conversations, err := h.conversationService.ListConversations(c.Context(), userID, cursorPtr, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve conversations", err)
	}

	return utils.SuccessResponse(c, "Conversations retrieved successfully", conversations)
}

// GetUnreadCount retrieves total unread message count
// GET /conversations/unread-count
func (h *ConversationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	unreadCount, err := h.conversationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve unread count", err)
	}

	return utils.SuccessResponse(c, "Unread count retrieved successfully", unreadCount)
}

// MarkAsRead marks all messages in a conversation as read
// POST /conversations/:conversationId/read
func (h *ConversationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Mark as read in database
	if err := h.conversationService.MarkAsRead(c.Context(), conversationID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to mark as read", err)
	}

	// Send WebSocket notification to the other user (sender)
	h.sendReadNotification(c, conversationID, userID)

	return utils.SuccessResponse(c, "Conversation marked as read", nil)
}

// sendReadNotification sends WebSocket notification to sender when receiver reads messages
func (h *ConversationHandler) sendReadNotification(c *fiber.Ctx, conversationID uuid.UUID, readerID uuid.UUID) {
	if h.chatHub == nil {
		log.Printf("⚠️ ChatHub is nil, skipping read notification")
		return
	}

	// Get conversation to find the other participant (sender)
	conversation, err := h.conversationRepo.GetByID(c.Context(), conversationID)
	if err != nil {
		log.Printf("⚠️ Failed to get conversation for read notification: %v", err)
		return
	}

	// Determine the other user (sender who should receive the notification)
	var senderID uuid.UUID
	if conversation.User1ID == readerID {
		senderID = conversation.User2ID
	} else {
		senderID = conversation.User1ID
	}

	// Send read update notification to sender
	h.chatHub.SendToUser(senderID, &chatWebsocket.ChatMessage{
		Type: "message.read_update",
		Payload: map[string]interface{}{
			"conversationId": conversationID.String(),
			"readBy":         readerID.String(),
			"readAt":         time.Now().Format(time.RFC3339),
		},
	})

	log.Printf("📤 Read notification sent to sender: %s (conversation: %s, read by: %s)", senderID, conversationID, readerID)
}

// SearchUsersForChat searches users for starting a new chat
// GET /chat/search-users?q=username&limit=20
func (h *ConversationHandler) SearchUsersForChat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get query params
	query := c.Query("q", "")
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}
	}

	// Search users
	results, err := h.conversationService.SearchUsersForChat(c.Context(), userID, query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search users", err)
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", results)
}
