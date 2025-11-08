package websocket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

// routeMessage routes incoming messages to appropriate handlers
func (h *ChatHub) routeMessage(client *ChatClient, message *ChatMessage) {
	ctx := context.Background()

	switch message.Type {
	// Message events
	case "message.send":
		h.handleMessageSend(ctx, client, message)

	case "message.read":
		h.handleMessageRead(ctx, client, message)

	// Typing indicators
	case "typing.start":
		h.handleTypingStart(ctx, client, message)

	case "typing.stop":
		h.handleTypingStop(ctx, client, message)

	// Online status (ping)
	case "ping":
		h.handlePing(ctx, client, message)

	// Block events
	case "block.add":
		h.handleBlockAdd(ctx, client, message)

	case "block.remove":
		h.handleBlockRemove(ctx, client, message)

	default:
		log.Printf("Unknown message type: %s", message.Type)
		client.sendError("unknown_type", "Unknown message type: "+message.Type)
	}
}

// ==================== Message Events ====================

// handleMessageSend handles incoming message from sender
func (h *ChatHub) handleMessageSend(ctx context.Context, client *ChatClient, message *ChatMessage) {
	// Parse conversationId
	conversationIDStr, ok := message.Payload["conversationId"].(string)
	if !ok {
		client.sendError("validation_error", "conversationId is required")
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		client.sendError("validation_error", "Invalid conversationId")
		return
	}

	// Parse content (optional if media is provided)
	var content *string
	if contentStr, ok := message.Payload["content"].(string); ok && contentStr != "" {
		content = &contentStr
	}

	// Parse type (default: text)
	messageType := "text"
	if typeStr, ok := message.Payload["type"].(string); ok {
		messageType = typeStr
	}

	// Parse media (optional)
	var media []dto.MessageMedia
	if mediaData, ok := message.Payload["media"].([]interface{}); ok && len(mediaData) > 0 {
		for _, m := range mediaData {
			if mediaMap, ok := m.(map[string]interface{}); ok {
				mediaItem := dto.MessageMedia{
					Type: mediaMap["type"].(string),
					URL:  mediaMap["url"].(string),
				}
				if size, ok := mediaMap["size"].(float64); ok {
					sizeVal := int64(size)
					mediaItem.Size = &sizeVal
				}
				if width, ok := mediaMap["width"].(float64); ok {
					w := int(width)
					mediaItem.Width = &w
				}
				if height, ok := mediaMap["height"].(float64); ok {
					h := int(height)
					mediaItem.Height = &h
				}
				if duration, ok := mediaMap["duration"].(float64); ok {
					d := int(duration)
					mediaItem.Duration = &d
				}
				media = append(media, mediaItem)
			}
		}
	}

	// Parse tempId (for client-side optimistic updates)
	var tempID *string
	if tempIDStr, ok := message.Payload["tempId"].(string); ok {
		tempID = &tempIDStr
	}

	// Validate: content OR media must be provided
	if (content == nil || *content == "") && len(media) == 0 {
		client.sendError("validation_error", "Either content or media must be provided")
		return
	}

	// Create SendMessageRequest
	req := &dto.SendMessageRequest{
		ConversationID: conversationID,
		Type:           messageType,
		Content:        content,
		Media:          media,
		TempID:         tempID,
	}

	// Send message via MessageService
	msgResponse, err := h.messageService.SendMessage(ctx, client.UserID, req)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		client.sendError("send_failed", err.Error())
		return
	}

	// Send acknowledgment to sender (message.sent)
	h.sendToClient(client, &ChatMessage{
		Type: "message.sent",
		Payload: map[string]interface{}{
			"message": msgResponse,
			"tempId":  tempID, // Echo back tempId for client matching
		},
	})

	// Send new message notification to receiver (message.new)
	receiverID := msgResponse.Receiver.ID
	h.sendToUser(receiverID, &ChatMessage{
		Type: "message.new",
		Payload: map[string]interface{}{
			"message": msgResponse,
		},
	})

	// If receiver is offline, send push notification
	if !h.IsUserOnline(receiverID) {
		go h.sendPushNotification(ctx, receiverID, client.UserID, msgResponse)
	}
}

// handleMessageRead handles mark as read request
func (h *ChatHub) handleMessageRead(ctx context.Context, client *ChatClient, message *ChatMessage) {
	// Parse conversationId
	conversationIDStr, ok := message.Payload["conversationId"].(string)
	if !ok {
		client.sendError("validation_error", "conversationId is required")
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		client.sendError("validation_error", "Invalid conversationId")
		return
	}

	// Mark all messages as read
	if err := h.conversationService.MarkAsRead(ctx, conversationID, client.UserID); err != nil {
		log.Printf("Failed to mark as read: %v", err)
		client.sendError("mark_read_failed", err.Error())
		return
	}

	// Send acknowledgment to reader (message.read_ack)
	h.sendToClient(client, &ChatMessage{
		Type: "message.read_ack",
		Payload: map[string]interface{}{
			"conversationId": conversationID.String(),
			"readAt":         time.Now().Format(time.RFC3339),
		},
	})

	// Get conversation to find other participant
	conversation, err := h.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get conversation: %v", err)
		return
	}

	// Determine other participant
	var otherUserID uuid.UUID
	if conversation.User1ID == client.UserID {
		otherUserID = conversation.User2ID
	} else {
		otherUserID = conversation.User1ID
	}

	// Send read update to sender (message.read_update)
	h.sendToUser(otherUserID, &ChatMessage{
		Type: "message.read_update",
		Payload: map[string]interface{}{
			"conversationId": conversationID.String(),
			"readBy":         client.UserID.String(),
			"readAt":         time.Now().Format(time.RFC3339),
		},
	})
}

// ==================== Typing Indicators ====================

// handleTypingStart broadcasts typing indicator
func (h *ChatHub) handleTypingStart(ctx context.Context, client *ChatClient, message *ChatMessage) {
	conversationIDStr, ok := message.Payload["conversationId"].(string)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		return
	}

	// Get conversation to find other participant
	conversation, err := h.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return
	}

	// Determine other participant
	var otherUserID uuid.UUID
	if conversation.User1ID == client.UserID {
		otherUserID = conversation.User2ID
	} else {
		otherUserID = conversation.User1ID
	}

	// Broadcast typing indicator to other user
	h.sendToUser(otherUserID, &ChatMessage{
		Type: "typing.start",
		Payload: map[string]interface{}{
			"conversationId": conversationID.String(),
			"userId":         client.UserID.String(),
		},
	})
}

// handleTypingStop broadcasts stop typing indicator
func (h *ChatHub) handleTypingStop(ctx context.Context, client *ChatClient, message *ChatMessage) {
	conversationIDStr, ok := message.Payload["conversationId"].(string)
	if !ok {
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		return
	}

	// Get conversation to find other participant
	conversation, err := h.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return
	}

	// Determine other participant
	var otherUserID uuid.UUID
	if conversation.User1ID == client.UserID {
		otherUserID = conversation.User2ID
	} else {
		otherUserID = conversation.User1ID
	}

	// Broadcast stop typing to other user
	h.sendToUser(otherUserID, &ChatMessage{
		Type: "typing.stop",
		Payload: map[string]interface{}{
			"conversationId": conversationID.String(),
			"userId":         client.UserID.String(),
		},
	})
}

// ==================== Ping/Pong ====================

// handlePing responds with pong and updates online status
func (h *ChatHub) handlePing(ctx context.Context, client *ChatClient, message *ChatMessage) {
	// Update online status in Redis
	_ = h.redisService.SetUserOnline(ctx, client.UserID)

	// Send pong response
	h.sendToClient(client, &ChatMessage{
		Type: "pong",
		Payload: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	})
}

// ==================== Block Events ====================

// handleBlockAdd handles block user request
func (h *ChatHub) handleBlockAdd(ctx context.Context, client *ChatClient, message *ChatMessage) {
	usernameToBlock, ok := message.Payload["username"].(string)
	if !ok {
		client.sendError("validation_error", "username is required")
		return
	}

	// Block user via BlockService
	if err := h.blockService.BlockUser(ctx, client.UserID, usernameToBlock); err != nil {
		log.Printf("Failed to block user: %v", err)
		client.sendError("block_failed", err.Error())
		return
	}

	// Send acknowledgment
	h.sendToClient(client, &ChatMessage{
		Type: "block.added",
		Payload: map[string]interface{}{
			"username":  usernameToBlock,
			"blockedAt": time.Now().Format(time.RFC3339),
		},
	})
}

// handleBlockRemove handles unblock user request
func (h *ChatHub) handleBlockRemove(ctx context.Context, client *ChatClient, message *ChatMessage) {
	usernameToUnblock, ok := message.Payload["username"].(string)
	if !ok {
		client.sendError("validation_error", "username is required")
		return
	}

	// Unblock user via BlockService
	if err := h.blockService.UnblockUser(ctx, client.UserID, usernameToUnblock); err != nil {
		log.Printf("Failed to unblock user: %v", err)
		client.sendError("unblock_failed", err.Error())
		return
	}

	// Send acknowledgment
	h.sendToClient(client, &ChatMessage{
		Type: "block.removed",
		Payload: map[string]interface{}{
			"username":    usernameToUnblock,
			"unblockedAt": time.Now().Format(time.RFC3339),
		},
	})
}

// ==================== Push Notifications ====================

// sendPushNotification sends push notification to offline user
func (h *ChatHub) sendPushNotification(ctx context.Context, receiverID uuid.UUID, senderID uuid.UUID, message *dto.MessageResponse) {
	// Get sender info for notification
	sender := message.Sender

	// Prepare notification title and body
	title := fmt.Sprintf("New message from %s", sender.Username)
	body := ""

	// Format body based on message type
	switch message.Type {
	case "text":
		if message.Content != nil {
			body = *message.Content
			// Truncate if too long
			if len(body) > 100 {
				body = body[:97] + "..."
			}
		} else {
			body = "Sent a message"
		}
	case "image":
		body = "📷 Sent a photo"
	case "video":
		body = "🎥 Sent a video"
	case "file":
		body = "📎 Sent a file"
	default:
		body = "Sent a message"
	}

	// Send push notification
	payload := &dto.PushNotificationPayload{
		Title: title,
		Body:  body,
		Icon:  "/icon-192x192.png",
		Tag:   "chat-message",
		Data: map[string]interface{}{
			"type":           "chat.message",
			"conversationId": message.ConversationID.String(),
			"messageId":      message.ID.String(),
			"senderId":       senderID.String(),
		},
	}

	err := h.pushService.SendToUser(ctx, receiverID, payload)

	if err != nil {
		log.Printf("Failed to send push notification to %s: %v", receiverID, err)
	} else {
		log.Printf("📬 Push notification sent to offline user %s", receiverID)
	}
}
