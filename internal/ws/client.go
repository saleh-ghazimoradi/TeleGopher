package ws

import (
	"github.com/coder/websocket"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"sync"
)

type Client struct {
	User *domain.User    `json:"user"`
	Conn *websocket.Conn `json:"-"`
	Send chan Event      `json:"-"`
	once sync.Once
}

func (c *Client) SendEvent(event Event) {
	select {
	case c.Send <- event:
	default:
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		if c.Conn != nil {
			_ = c.Conn.Close(websocket.StatusNormalClosure, "Closing connection")
		}
		close(c.Send)
	})
}

func NewClient(user *domain.User, conn *websocket.Conn) *Client {
	return &Client{
		User: user,
		Conn: conn,
		Send: make(chan Event),
	}
}
