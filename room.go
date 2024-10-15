package main

import (
	"strconv"
)

type RoomState int

const (
	WaitingRoomState  RoomState = 0
	RunningRoomState            = 1
	FinishedRoomState           = 2
)

type Room struct {
	roomState RoomState
	game      *Game

	register   chan *Client
	unregister chan *Client

	clients map[*Client]bool

	broadcast chan []byte
}

func newRoom() *Room {
	return &Room{
		roomState:  WaitingRoomState,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
	}
}

func (r *Room) broadcastBytes(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}

func (r *Room) broadcastString(message string) {
	r.broadcastBytes([]byte(message))
}

func (r *Room) run() {
	actions := make(chan Action)

	r.game = newGame(actions, r)
	go r.game.run()
	pid := 0

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			p := newPlayer("ogi"+strconv.FormatInt(int64(pid), 10), client)
			r.game.addPlayer(p)
			pid++
			client.player = p
			r.broadcastString(p.name + "joined!")
		case client := <-r.unregister:
			r.clients[client] = false
			r.game.removePlayer(client.player)
		case message := <-r.broadcast:
			r.broadcastBytes(message)
		}
	}
}
