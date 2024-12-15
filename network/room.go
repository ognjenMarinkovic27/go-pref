package network

import (
	"maps"
	"ognjen/go-pref/game"
	"slices"
)

type RoomState int

const (
	WaitingRoomState  RoomState = 0
	RunningRoomState  RoomState = 1
	FinishedRoomState RoomState = 2
)

type Room struct {
	roomState RoomState
	game      *game.Game

	register   chan *Client
	unregister chan *Client

	clients map[string]*Client

	recv chan InboundMessage
	seq  int
}

func NewRoom() *Room {
	return &Room{
		roomState:  WaitingRoomState,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		recv:       make(chan InboundMessage),
		seq:        0,
	}
}

func (r *Room) handleAction(action game.Action) bool {
	g := r.game
	if g.Validate(action) {
		g.Apply(action)

		responses := g.Collect()
		messages := r.messagesFromResponses(responses)
		for _, m := range messages {
			r.sendMessage(m)
		}

		return true
	}

	return false
}

func (room *Room) messagesFromResponses(responses []game.Response) []OutboundMessage {
	var messages []OutboundMessage
	for _, r := range responses {
		messages = append(messages, room.messageFromResponse(r))
	}

	return messages
}

func (r *Room) messageFromResponse(response game.Response) OutboundMessage {
	client := r.clientFromPid(response.RecepientPid())
	return OutboundMessage{
		Recepient: client,
		Payload:   response,
	}
}

func (r *Room) clientFromPid(pid string) *Client {
	var client *Client
	if pid == "" {
		client = nil
	} else {
		client = r.clients[pid]
	}

	return client
}

func (r *Room) sendMessage(m OutboundMessage) {
	if m.Recepient == nil {
		r.broadcastMessage(m)
		return
	}
	m.Recepient.send <- m
}

func (r *Room) broadcastMessage(m OutboundMessage) {
	for _, c := range r.clients {
		c.send <- m
	}
}

func (r *Room) Run() {
	r.game = game.NewGame()

	for {
		select {
		case client := <-r.register:
			r.broadcastMessage(r.NewServerMessage(nil, &ConnectedPayload{client.pid}))
			if r.clients[client.pid] == nil {
				r.game.AddPlayer(client.pid)
			}
			r.clients[client.pid] = client
			r.sendLobby(client)
		case client := <-r.unregister:
			if !r.game.Started() {
				r.game.RemovePlayer(client.pid)
			}
			delete(r.clients, client.pid)
			r.broadcastMessage(r.NewServerMessage(nil, &DisconnectedPayload{client.pid}))
		case message := <-r.recv:
			action := message.Payload.(game.Action)
			action.SetPlayerPid(message.Client.pid)
			if !r.handleAction(action) {
				r.sendMessage(OutboundMessage{
					Seq:       message.Seq,
					Recepient: message.Client,
					Payload:   &InvalidAction{},
				})
			}
		}
	}
}

func (r *Room) sendLobby(rec *Client) {
	pids := slices.Collect(maps.Keys(r.clients))
	r.sendMessage(r.NewServerMessage(rec, &LobbyPayload{pids}))
}

func (r *Room) NewServerMessage(recepient *Client, payload Payload) OutboundMessage {
	msg := OutboundMessage{
		Seq:       r.seq,
		Recepient: recepient,
		Payload:   payload,
	}

	r.seq++

	return msg
}
