package websocket

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/command"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "websocket")

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

// Client is a connected websocket client.
type Client struct {
	uuid         string
	hub          *Hub
	conn         *websocket.Conn
	send         chan []byte
	done         chan struct{}
	commandQueue chan<- command.Command
}

// NewClient creates a new *Client to handle a websocket connection.
func NewClient(hub *Hub, conn *websocket.Conn, queue chan<- command.Command) *Client {
	return &Client{
		uuid:         uuid.NewV4().String(),
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, sendBufSize),
		done:         make(chan struct{}),
		commandQueue: queue,
	}
}

// Listen registers the client to the websocket hub and starts listening for
// incoming data from and data that should be written to the websocket.
func (c *Client) Listen() {
	c.hub.register <- c

	go c.listenWrite()
	go c.listenRead()
}

// listenRead reads messages from the websocket and processes them.
func (c *Client) listenRead() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case <-c.done:
			return
		default:
			envelope := Envelope{}

			if err := c.conn.ReadJSON(&envelope); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Errorf("websocket read error: %v", err)
				}
				return
			}

			cmd, err := decodeCommand(envelope)
			if err != nil {
				log.Errorf("failed to decode command: %v", err)
				continue
			}

			if clientAwareCmd, ok := cmd.(command.SenderAwareCommand); ok {
				clientAwareCmd.SetSender(c)
			}

			c.commandQueue <- cmd
		}
	}
}

// listenWrite writes messages received from the hub back to the websocket
// connection.
func (c *Client) listenWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.done:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.WithFields(logrus.Fields{
				"length": len(message),
				"uuid":   c.uuid,
			}).Debug("sending message")

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Errorf("websocket write error: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("websocket write error: %v", err)
				return
			}
		}
	}
}

// Send sends msg through the websocket to the connected client.
func (c *Client) Send(msg []byte) {
	c.send <- msg
}
