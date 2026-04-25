package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/ws"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebSocketHandler struct {
	userService    service.UserService
	messageService service.MessageService
	logger         utils.LoggerStrategy
	hub            *ws.Hub
	cfg            *config.Config
}

func (wsh *WebSocketHandler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		helper.UnauthorizedResponse(w, "Authorization header missing")
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		helper.UnauthorizedResponse(w, "Invalid authorization header format")
		return
	}

	claims, err := utils.ValidateToken(tokenParts[1], wsh.cfg.JWT.Secret)
	if err != nil {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	user, err := wsh.userService.GetUserById(r.Context(), claims.UserId)
	if err != nil {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		helper.InternalServerError(w, "failed to accept websocket connection", err)
		return
	}

	client := ws.NewClient(&domain.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, conn)

	wsh.hub.RegisterClient(client)
	wsh.hub.SendCurrentClients(client)

	defer func() {
		wsh.hub.UnregisterClient(client)
		client.Close()
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go wsh.heartbeat(ctx, client)
	go wsh.writePump(ctx, client)
	wsh.readPump(ctx, cancel, client)
}

func (wsh *WebSocketHandler) heartbeat(ctx context.Context, client *ws.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := client.Conn.Ping(pingCtx)
			if err != nil {
				wsh.logger.Warn("Ping failed for client", "client", client.User.Id, "err", err)
				cancel()
				client.Close()
				return
			}
			cancel()

			client.SendEvent(ws.Event{
				EventType: ws.EventHeartbeat,
				Payload:   nil,
			})
		}
	}
}

func (wsh *WebSocketHandler) writePump(ctx context.Context, client *ws.Client) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-client.Send:
			if !ok {
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := client.Conn.Write(writeCtx, websocket.MessageText, wsh.eventToJSON(event))
			cancel()

			if err != nil {
				wsh.logger.Error("failed to write to client", "client", client.User.Id, "err", err)
				return
			}
		}
	}
}

func (wsh *WebSocketHandler) readPump(ctx context.Context, cancel context.CancelFunc, client *ws.Client) {
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			wsh.logger.Error("Recovered from panic in readPump for client", "client", client.User.Id, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, message, err := client.Conn.Read(ctx)
		if err != nil {
			return
		}

		var event ws.Event
		if err := json.Unmarshal(message, &event); err != nil {
			wsh.hub.SendError(client.User.Id, "invalid event format")
			continue
		}

		wsh.handleIncomingEvent(client, event)
	}
}

func (wsh *WebSocketHandler) handleIncomingEvent(client *ws.Client, event ws.Event) {
	payload, ok := event.Payload.(map[string]any)
	if !ok {
		wsh.hub.SendError(client.User.Id, "invalid event format")
		return
	}

	switch event.EventType {
	case ws.EventMessage:
		wsh.handleMessageEvent(client, payload)
	case ws.EventDelivered:
		wsh.handleDeliveredEvent(client, payload)
	case ws.EventRead:
		wsh.handleReadEvent(client, payload)
	case ws.EventTyping:
		wsh.handleTypingEvent(client, payload)
	default:
		wsh.hub.SendError(client.User.Id, "unknown event type: "+string(event.EventType))
	}
}

func (wsh *WebSocketHandler) handleMessageEvent(client *ws.Client, payload map[string]any) {
	privateId, ok := wsh.extractUint(payload, "private_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "private_id is required and must be a number")
		return
	}

	receiverId, ok := wsh.extractUint(payload, "receiver_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "receiver_id is required and must be a number")
		return
	}

	messageType, ok := payload["message_type"].(string)
	if !ok || messageType == "" {
		wsh.hub.SendError(client.User.Id, "message_type is required")
		return
	}

	content, ok := payload["content"].(string)
	if !ok || content == "" {
		wsh.hub.SendError(client.User.Id, "content is required")
		return
	}

	// Create message via service
	req := &dto.MessageRequest{
		PrivateId:   privateId,
		MessageType: messageType,
		Content:     content,
	}

	message, err := wsh.messageService.SendMessage(context.Background(), req, client.User.Id)
	if err != nil {
		wsh.hub.SendError(client.User.Id, fmt.Sprintf("failed to send message: %v", err))
		return
	}

	// Broadcast to sender and receiver
	wsh.hub.SendEventToUserIds([]uint{client.User.Id, receiverId}, client.User.Id, ws.EventMessage, map[string]any{
		"message": message,
	})
}

func (wsh *WebSocketHandler) handleDeliveredEvent(client *ws.Client, payload map[string]any) {
	messageId, ok := wsh.extractUint(payload, "message_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "message_id is required and must be a number")
		return
	}

	err := wsh.messageService.MarkMessageAsDelivered(context.Background(), messageId, client.User.Id)
	if err != nil {
		wsh.hub.SendError(client.User.Id, fmt.Sprintf("failed to mark message as delivered: %v", err))
		return
	}

	// Get message to notify sender
	message, err := wsh.messageService.GetMessage(context.Background(), messageId, client.User.Id)
	if err != nil {
		// Still successful for the recipient, just log error
		log.Printf("Failed to get message for notification: %v", err)
		return
	}

	// Notify sender that message was delivered
	wsh.hub.SendEventToUserIds([]uint{message.FromId}, client.User.Id, ws.EventDelivered, map[string]any{
		"message_id": messageId,
		"to_id":      client.User.Id,
	})
}

func (wsh *WebSocketHandler) handleReadEvent(client *ws.Client, payload map[string]any) {
	messageId, ok := wsh.extractUint(payload, "message_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "message_id is required and must be a number")
		return
	}

	err := wsh.messageService.MarkMessageAsRead(context.Background(), messageId, client.User.Id)
	if err != nil {
		wsh.hub.SendError(client.User.Id, fmt.Sprintf("failed to mark message as read: %v", err))
		return
	}

	// Get message to notify sender
	message, err := wsh.messageService.GetMessage(context.Background(), messageId, client.User.Id)
	if err != nil {
		log.Printf("Failed to get message for notification: %v", err)
		return
	}

	// Notify sender that message was read
	wsh.hub.SendEventToUserIds([]uint{message.FromId}, client.User.Id, ws.EventRead, map[string]any{
		"message_id": messageId,
	})
}

func (wsh *WebSocketHandler) handleTypingEvent(client *ws.Client, payload map[string]any) {
	privateId, ok := wsh.extractUint(payload, "private_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "private_id is required and must be a number")
		return
	}

	receiverId, ok := wsh.extractUint(payload, "receiver_id")
	if !ok {
		wsh.hub.SendError(client.User.Id, "receiver_id is required and must be a number")
		return
	}

	isTyping, ok := payload["is_typing"].(bool)
	if !ok {
		wsh.hub.SendError(client.User.Id, "is_typing is required and must be a boolean")
		return
	}

	wsh.hub.SendEventToUserIds([]uint{receiverId}, client.User.Id, ws.EventTyping, map[string]any{
		"private_id": privateId,
		"user_id":    client.User.Id,
		"is_typing":  isTyping,
	})
}

func (wsh *WebSocketHandler) extractUint(payload map[string]any, key string) (uint, bool) {
	value, ok := payload[key]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case float64:
		return uint(v), true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	default:
		return 0, false
	}
}

func (wsh *WebSocketHandler) eventToJSON(event ws.Event) []byte {
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return []byte(`{"event_type":"error","payload":"internal server error"}`)
	}
	return jsonData
}

func NewWebSocketHandler(userService service.UserService, messageService service.MessageService, logger utils.LoggerStrategy, hub *ws.Hub, cfg *config.Config) *WebSocketHandler {
	return &WebSocketHandler{
		userService:    userService,
		messageService: messageService,
		logger:         logger,
		hub:            hub,
		cfg:            cfg,
	}
}
