package control

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/message"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	sendBufSize = 256
)

// Client type definition
type Client struct {
	hub        *Hub
	dispatcher message.Dispatcher
	conn       *websocket.Conn
	send       chan []byte
}

// NewClient create a new client to handle a websocket connection
func NewClient(hub *Hub, dispatcher message.Dispatcher, conn *websocket.Conn) *Client {
	return &Client{
		hub:        hub,
		dispatcher: dispatcher,
		conn:       conn,
		send:       make(chan []byte, sendBufSize),
	}
}

func (c *Client) Listen() {
	c.hub.register <- c

	go c.listenWrite()
	go c.listenRead()
}

// listenRead reads messages from the websocket and processes them
func (c *Client) listenRead() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		msg := message.Envelope{}

		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}

		if err := c.dispatcher.Dispatch(msg); err != nil {
			log.Println(err)
		}
	}
}

// listenWrite writes messages received from the hub back to the websocket connection
func (c *Client) listenWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
