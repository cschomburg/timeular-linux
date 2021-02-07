package timeular

type Client struct {
	hub  *Hub
	send chan Timeular
}

type Hub struct {
	clients   map[*Client]bool
	broadcast chan Timeular
	register  chan *Client
	lastState Timeular
}

func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan Timeular),
		register:  make(chan *Client),
		clients:   make(map[*Client]bool),
	}
}

func (h *Hub) Register(send chan Timeular) {
	c := &Client{send: send}
	h.register <- c
}

func (h *Hub) Broadcast(t Timeular) {
	h.broadcast <- t
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if !h.lastState.Tracking.StartedAt.IsZero() {
				client.send <- h.lastState
			}

		case message := <-h.broadcast:
			h.lastState = message
			for client := range h.clients {
				select {
				case client.send <- message:
				}
			}
		}
	}
}
