package main

import (
	"fmt"
)

type PlayerScore struct {
	prevSoups int
	score     int
	nextSoups int
}

type Player struct {
	name  string
	id    int
	hand  [10]Card
	score PlayerScore

	client *Client
}

func newPlayer(name string, id int, client *Client) *Player {
	return &Player{
		name: name,
		id:   id,
		score: PlayerScore{
			prevSoups: 0,
			score:     0,
			nextSoups: 0,
		},
		client: client,
	}
}

func messageToAction(message []byte, p *Player) Action {
	pi := PlayerInfo{player: p}

	strMessage := string(message)

	switch strMessage[0] {
	case 'r':
		return ReadyAction{PlayerInfo: pi}
	case 'b':
		return BidAction{PlayerInfo: pi}
	case 'p':
		return PassBidAction{PlayerInfo: pi}
	default:
		fmt.Println("Invalid action from", p.name)
		return nil
	}
}
