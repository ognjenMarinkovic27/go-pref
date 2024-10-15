package main

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

	recv chan []byte
	sender *Client
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

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			p := newPlayer(client)
			r.game.addPlayer(p)
			client.player = p
			r.broadcastString(p.getName() + "joined!")
		case client := <-r.unregister:
			r.clients[client] = false
			r.game.removePlayer(client.player)
		case message := <-r.broadcast:
			r.broadcastBytes(message)
		case message := <-r.recv:
			action := messageToAction(message, r.sender.player)
			r.game.actions <- action
		}
	}
}
