package websocket

// Hub acts as a central registry for connected websocket clients and can be
// used to broadcast messages to everyone.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// NewHub creates a new hub for handling communicating between connected
// websocket clients.
func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}

	return h
}

// Run runs the control loop. If stopCh is closed, the hub will disconnect all
// clients and stop the control loop.
func (h *Hub) Run(stopCh <-chan struct{}) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Infof("registered new client %s", client.uuid)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.done)
				delete(h.clients, client)
				log.Infof("unregistered client %s", client.uuid)
			}
		case msg := <-h.broadcast:
			log.WithField("length", len(msg)).Debug("broadcasting message")
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					h.unregister <- client
				}
			}
		case <-stopCh:
			log.Infof("shutting down hub")
			for client := range h.clients {
				close(client.done)
				delete(h.clients, client)
			}
			return
		}
	}
}

// Broadcast broadcasts msg to all connected clients.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}
