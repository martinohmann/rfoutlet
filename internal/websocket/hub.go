// Package websocket provides the http handler, client and hub for managing the
// websocket connections of clients using rfoutlet. This is both used to react
// on commands of a single client as well as to broadcast updates to all
// clients so that state changes are immediately visible to everybody.
package websocket

// Hub acts as a central registry for connected websocket clients and can be
// used to broadcast messages to everyone.
type Hub struct {
	clients    map[*client]struct{}
	register   chan *client
	unregister chan *client
	broadcast  chan []byte
	send       chan clientMsg
}

// NewHub creates a new hub for handling communicating between connected
// websocket clients.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*client]struct{}),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan []byte),
		send:       make(chan clientMsg),
	}
}

// Run runs the control loop. If stopCh is closed, the hub will disconnect all
// clients and stop the control loop.
func (h *Hub) Run(stopCh <-chan struct{}) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			log.WithField("uuid", client.uuid).Info("new client registered")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.unregisterClient(client)
			}
		case cm := <-h.send:
			if _, ok := h.clients[cm.client]; ok {
				h.sendClientMessage(cm.client, cm.msg)
			}
		case msg := <-h.broadcast:
			log.WithField("length", len(msg)).Debug("broadcasting message")
			for client := range h.clients {
				h.sendClientMessage(client, msg)
			}
		case <-stopCh:
			log.Infof("shutting down hub")
			for client := range h.clients {
				h.unregisterClient(client)
			}
			return
		}
	}
}

// Broadcast broadcasts msg to all connected clients.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

// Send sends a message to a specific client.
func (h *Hub) Send(client *client, msg []byte) {
	h.send <- clientMsg{client, msg}
}

func (h *Hub) unregisterClient(client *client) {
	close(client.send)
	delete(h.clients, client)
	log.WithField("uuid", client.uuid).Info("client unregistered")
}

func (h *Hub) sendClientMessage(client *client, msg []byte) {
	select {
	case client.send <- msg:
	default:
		h.unregister <- client
	}
}

// clientMsg is a wrapper type for a message destined for a specific client.
type clientMsg struct {
	client *client
	msg    []byte
}
