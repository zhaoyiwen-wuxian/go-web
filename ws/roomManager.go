package ws

import (
	"encoding/json"
	"go-web/models"
)

func (rm *RoomManager) JoinRoom(roomID, userID string, client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.Rooms[roomID]
	if !exists {
		room = &Room{ID: roomID, Clients: make(map[string]*Client)}
		rm.Rooms[roomID] = room
	}
	room.Clients[userID] = client
}

func (rm *RoomManager) LeaveRoom(roomID, userID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, ok := rm.Rooms[roomID]; ok {
		delete(room.Clients, userID)
		if len(room.Clients) == 0 {
			delete(rm.Rooms, roomID)
		}
	}
}

func (rm *RoomManager) Broadcast(roomID, from string, msg models.Message) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if room, ok := rm.Rooms[roomID]; ok {
		for uid, client := range room.Clients {
			if uid != from {
				client.Send <- marshal(msg)
			}
		}
	}
}

func marshal(msg models.Message) []byte {
	b, _ := json.Marshal(msg)
	return b
}
