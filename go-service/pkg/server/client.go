package server

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	id     string
	ctx    context.Context
	cancel context.CancelFunc
	inMsg  chan []byte
}

func NewClient(conn *websocket.Conn, id string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		conn:   conn,
		id:     id,
		ctx:    ctx,
		cancel: cancel,
		inMsg:  make(chan []byte),
	}
}
func (c *Client) GetId() string {
	return c.id
}

func (c *Client) Close() {
	close(c.inMsg)
	c.conn.Close()
	c.cancel()
}

func (c *Client) GetMessageChann() <-chan []byte {

	out := make(chan []byte)

	go func() {
		defer close(out)
		defer c.Close()

		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					return
				}
				out <- message
			}
		}

	}()
	return out
}
