package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupChatRoutes(api fiber.Router, h *handlers.Handlers) {
	// All chat routes require authentication
	chat := api.Group("/chat", middleware.Protected())

	// Search users for chat
	chat.Get("/search-users", h.ConversationHandler.SearchUsersForChat)

	// Conversation routes (Nested URLs)
	conversations := chat.Group("/conversations")
	conversations.Get("/with/:username", h.ConversationHandler.GetOrCreateConversation)
	conversations.Get("/", h.ConversationHandler.ListConversations)
	conversations.Get("/unread-count", h.ConversationHandler.GetUnreadCount)

	// Nested message routes under conversations
	conversations.Get("/:conversationId/messages", h.MessageHandler.ListMessages)
	conversations.Post("/:conversationId/messages", h.MessageHandler.SendMessage)
	conversations.Post("/:conversationId/read", h.ConversationHandler.MarkAsRead)

	// Phase 2: Media/Links/Files filtering
	conversations.Get("/:conversationId/media", h.MessageHandler.GetConversationMedia)
	conversations.Get("/:conversationId/links", h.MessageHandler.GetConversationLinks)
	conversations.Get("/:conversationId/files", h.MessageHandler.GetConversationFiles)

	// Message routes (for message-specific operations)
	messages := chat.Group("/messages")
	messages.Get("/:id/context", h.MessageHandler.GetMessageContext)
	messages.Get("/:id", h.MessageHandler.GetMessage)

	// Block routes
	blocks := chat.Group("/blocks")
	blocks.Post("/", h.BlockHandler.BlockUser)
	blocks.Delete("/:username", h.BlockHandler.UnblockUser)
	blocks.Get("/status/:username", h.BlockHandler.GetBlockStatus)
	blocks.Get("/", h.BlockHandler.ListBlockedUsers)
}
