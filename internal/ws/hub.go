package ws

import (
	"context"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"log"
	"sync"
)

type Hub struct {
	privateService service.PrivateService
	messageService service.MessageService
	clients        map[int64]map[*Client]struct{}
	mu             sync.RWMutex
}

func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	conns, ok := h.clients[client.User.Id]
	if !ok {
		conns = make(map[*Client]struct{})
		h.clients[client.User.Id] = conns
	}
	conns[client] = struct{}{}
	firstConnection := len(conns) == 1
	h.mu.Unlock()

	if firstConnection {
		h.BroadcastToAll(Event{
			EventType: EventUserOnline,
			Payload:   client.User.ToMap(),
		})

		go func() {
			ctx := context.Background()
			privates, err := h.privateService.GetUserPrivates(ctx, client.User.Id)
			if err != nil {
				log.Println("failed to get privates:", err)
				return
			}
			for _, private := range privates {
				msgs, err := h.messageService.GetUndeliveredMessages(ctx, private.Id, client.User.Id)
				if err != nil {
					log.Println("failed to get undelivered messages:", err)
					continue
				}
				for _, msg := range msgs {
					if msg.FromId == client.User.Id {
						continue
					}
					h.SendEventToUserIds([]int64{msg.FromId}, client.User.Id, EventUserOnline, map[string]any{
						"message_id": msg.Id,
						"to_id":      client.User.Id,
					})
				}
			}
		}()
	}
}

func (h *Hub) SendCurrentClients(client *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]map[string]any, 0)
	seen := make(map[int64]struct{})

	for userId, conns := range h.clients {
		if userId == client.User.Id {
			continue
		}

		_, ok := seen[userId]
		if ok {
			continue
		}

		for c := range conns {
			users = append(users, c.User.ToMap())
			seen[userId] = struct{}{}
			break
		}
	}

	client.Send <- Event{
		EventType: EventCurrentUsers,
		Payload:   users,
	}
}

func (h *Hub) UnregisterClient(client *Client) {
	h.mu.Lock()
	conns, ok := h.clients[client.User.Id]
	if !ok {
		h.mu.Unlock()
		return
	}

	delete(conns, client)
	noConnectionLeft := len(conns) == 0
	if noConnectionLeft {
		delete(h.clients, client.User.Id)
	}

	h.mu.Unlock()

	if noConnectionLeft {
		h.BroadcastToAll(Event{
			EventType: EventUserOffline,
			Payload:   client.User.ToMap(),
		})
	}
}

func (h *Hub) GetClients(userId int64) ([]*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.clients[userId]
	if !ok || len(conns) == 0 {
		return nil, false
	}

	clients := make([]*Client, 0, len(conns))
	for conn := range conns {
		clients = append(clients, conn)
	}

	return clients, true
}

func (h *Hub) SendEventToUserIds(userIds []int64, senderId int64, event EventType, payload map[string]any) {
	for _, userId := range userIds {
		h.mu.RLock()
		conns, ok := h.clients[userId]
		h.mu.RUnlock()
		if !ok {
			continue
		}

		for conn := range conns {
			conn.SendEvent(Event{
				EventType: event,
				Payload:   payload,
			})
		}
	}
}

func (h *Hub) BroadcastToAll(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients {
		for c := range client {
			select {
			case c.Send <- event:
			default:
				log.Printf("warning: dropped event for client %d, channel full", c.User.Id)
			}
		}
	}
}

func (h *Hub) SendError(clientId int64, message string) {
	clients, ok := h.GetClients(clientId)
	if !ok || len(clients) == 0 {
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

	log.Println("shutting down Hub, notifying all clients...")

	for _, conns := range h.clients {
		for client := range conns {
			client.SendEvent(Event{
				EventType: EventServerShutdown,
				Payload:   "Server is shutting down",
			})
			client.Close()
		}
	}
	h.clients = make(map[int64]map[*Client]struct{})
	log.Println("Hub shutdown complete")
}

func NewHub(messageService service.MessageService, privateService service.PrivateService) *Hub {
	return &Hub{
		messageService: messageService,
		privateService: privateService,
		clients:        make(map[int64]map[*Client]struct{}),
	}
}
