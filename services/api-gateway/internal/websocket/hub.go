package websocket

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Connection]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*Connection]bool),
	}
}

func (h *Hub) Add(key string, c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[key] == nil {
		h.clients[key] = make(map[*Connection]bool)
	}
	h.clients[key][c] = true
}

func (h *Hub) Remove(key string, c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conns, ok := h.clients[key]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.clients, key)
		}
	}
}

func (h *Hub) SendTo(key string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conns, ok := h.clients[key]; ok {
		for c := range conns {
			_ = c.Send(data)
		}
	}
}