package ws

import (
	"sync"
)

type Room struct {
	ID      string
	Clients map[string]*Client
}

type RoomManager struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

var Manager = &RoomManager{
	Rooms: make(map[string]*Room),
}
