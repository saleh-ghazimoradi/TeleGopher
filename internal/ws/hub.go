package ws

import (
	"context"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"sync"
)

type Hub struct {
	Clients        map[uint]map[*Client]struct{}
	privateService service.PrivateService
	messageService service.MessageService
	logger         utils.LoggerStrategy
	mu             sync.RWMutex
}

func (h *Hub) BroadcastToAll(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.Clients {
		for client := range conn {
			select {
			case client.Send <- event:
			default:
				h.logger.Warn("dropped event for client", "client", client.User.Id, "channel full")
			}
		}
	}
}

func (h *Hub) GetClients(userId uint) ([]*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connections, ok := h.Clients[userId]
	if !ok || len(connections) == 0 {
		return nil, false
	}

	clients := make([]*Client, 0, len(connections))
	for c := range connections {
		clients = append(clients, c)
	}

	return clients, true
}

func (h *Hub) SendEventToUserIds(userIds []uint, senderId uint, eventType EventType, payload map[string]any) {
	for _, id := range userIds {
		h.mu.RLock()
		connections, ok := h.Clients[id]
		h.mu.RUnlock()
		if !ok {
			continue
		}

		for c := range connections {
			c.SendEvent(Event{
				EventType: eventType,
				Payload:   payload,
			})
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	connections, ok := h.Clients[client.User.Id]
	if !ok {
		connections = make(map[*Client]struct{})
		h.Clients[client.User.Id] = connections
	}

	connections[client] = struct{}{}
	firstConnection := len(connections) == 1
	h.mu.Unlock()

	if firstConnection {
		h.BroadcastToAll(Event{
			EventType: EventUserOnline,
			Payload:   client.User.ToMap(),
		})

		go func() {
			ctx := context.Background()
			privates, err := h.privateService.GetPrivatesForUser(ctx, client.User.Id)
			if err != nil {
				h.logger.Error("failed to get privates", "err", err)
				return
			}

			for _, private := range privates {
				msg, err := h.messageService.GetUndeliveredMessages(ctx, private.Id, client.User.Id)
				if err != nil {
					h.logger.Error("failed to get undelivered messages", "err", err)
					continue
				}
				for _, m := range msg {
					if m.FromId == client.User.Id {
						continue
					}
					h.SendEventToUserIds([]uint{m.FromId}, client.User.Id, EventUserOnline, map[string]any{
						"message_id": m.Id,
						"to_id":      client.User.Id,
					})
				}
			}
		}()
	}
}

func (h *Hub) UnregisterClient(client *Client) {
	h.mu.Lock()
	connections, ok := h.Clients[client.User.Id]
	if !ok {
		h.mu.Unlock()
		return
	}

	delete(connections, client)
	noConnectionLeft := len(connections) == 0
	if noConnectionLeft {
		delete(h.Clients, client.User.Id)
	}

	h.mu.Unlock()

	if noConnectionLeft {
		h.BroadcastToAll(Event{
			EventType: EventUserOffline,
			Payload:   client.User.ToMap(),
		})
	}
}

func (h *Hub) SendCurrentClients(client *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]map[string]any, 0)
	seen := make(map[uint]struct{})

	for userId, connections := range h.Clients {
		if userId == client.User.Id {
			continue
		}

		_, ok := seen[userId]
		if ok {
			continue
		}

		for connection := range connections {
			users = append(users, connection.User.ToMap())
			seen[userId] = struct{}{}
			break
		}
	}

	client.Send <- Event{
		EventType: EventCurrentUsers,
		Payload:   users,
	}
}

func (h *Hub) SendError(clientId uint, message string) {
	clients, ok := h.GetClients(clientId)
	if !ok || len(message) == 0 {
		return
	}

	for _, client := range clients {
		client.SendEvent(Event{
			EventType: EventError,
			Payload:   message,
		})
	}
}

func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.logger.Info("Shutting down hub, notifying all clients...")

	for _, connections := range h.Clients {
		for client := range connections {
			client.SendEvent(Event{
				EventType: EventServerShutdown,
				Payload:   "Server is shutting down",
			})
			client.Close()
		}
	}
	h.Clients = make(map[uint]map[*Client]struct{})
	h.logger.Info("Hub shutdown complete")
}

func NewHub(privateService service.PrivateService, messageService service.MessageService, logger utils.LoggerStrategy) *Hub {
	return &Hub{
		Clients:        make(map[uint]map[*Client]struct{}),
		privateService: privateService,
		messageService: messageService,
		logger:         logger,
	}
}
