package network

import (
	"ognjen/go-pref/game"
)

type RoomState int

const (
	WaitingRoomState  RoomState = 0
	RunningRoomState            = 1
	FinishedRoomState           = 2
)

type Room struct {
	roomState RoomState
	game      *game.Game

	register   chan *Client
	unregister chan *Client

	clients map[*Client]bool

	broadcast chan []byte

	recv chan InboundMessage
}

func NewRoom() *Room {
	return &Room{
		roomState:  WaitingRoomState,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		recv:       make(chan InboundMessage),
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

func (r *Room) handleAction(action game.Action) {
	g := r.game
	if g.Validate(action) {
		g.Apply(action)

		_ = g.Collect()
	}
}

func (r *Room) Run() {
	r.game = game.NewGame()

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			p := game.NewPlayer()
			r.game.AddPlayer(p)
			client.player = p
			// r.broadcastString(client.name + " joined!")
		case client := <-r.unregister:
			r.clients[client] = false
			r.game.RemovePlayer(client.player)
		case message := <-r.broadcast:
			r.broadcastBytes(message)
		case message := <-r.recv:
			action := message.Action()
			r.handleAction(action)
		}
	}
}
