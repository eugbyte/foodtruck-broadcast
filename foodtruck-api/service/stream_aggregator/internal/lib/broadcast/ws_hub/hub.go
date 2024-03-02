package wshub

import (
	ws "foodtruck/pkg/websocket/server"
	"sync"
)

type WSServer = ws.WSServer

type Hub struct {
	mu        sync.Mutex
	broadCast chan []byte
	wsMap     map[string]WSServer
}

// New hub to register websocket connections
func New() *Hub {
	return &Hub{
		broadCast: make(chan []byte),
		wsMap:     make(map[string]WSServer),
	}
}

func (hub *Hub) Conns() map[string]WSServer {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	conns := hub.wsMap
	return conns
}

func (hub *Hub) Add(wsID string, wsServer WSServer) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.wsMap[wsID] = wsServer
}

func (hub *Hub) Remove(wsID string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.wsMap, wsID)
}

func (hub *Hub) BroadcastCh() (broadcastChan chan []byte) {
	return hub.broadCast
}

// Upon receiving a message from the broadcast channel, further broadcast the same message to all registered connections.
// Blocking operation.
func (hub *Hub) Broadcast() {
	for msg := range hub.broadCast {
		for _, conn := range hub.Conns() {
			conn.Channel() <- msg
		}
	}
}
