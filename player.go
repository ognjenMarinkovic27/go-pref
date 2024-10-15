package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type PlayerScore struct {
	score int
	soups map[*Player]int
}

type Player struct {
	name   string
	hand   [10]Card
	played map[Card]bool
	score  PlayerScore

	client *Client
	next   *Player
}

func newPlayer(name string, client *Client) *Player {
	p := &Player{
		name: name,
		score: PlayerScore{
			soups: make(map[*Player]int),
		},
		played: make(map[Card]bool),
		client: client,
	}
	p.next = p
	return p
}

func messageToAction(message []byte, p *Player) Action {
	pi := PlayerInfo{player: p}

	strMessage := string(message)

	// TODO: clean up
	switch strMessage[0] {
	case 'r':
		return ReadyAction{PlayerInfo: pi}
	case 'b':
		return BidAction{PlayerInfo: pi}
	case 'p':
		return PassBidAction{PlayerInfo: pi}
	case 'c':
		if len(strMessage) > 1 {
			val, err := strconv.Atoi(strMessage[1:])
			if err == nil && val <= int(SansGameType) {
				return ChooseGameTypeAction{gameType: GameType(val), PlayerInfo: pi}
			}
		}
		return InvalidAction{PlayerInfo: pi}
	case 'n':
		return RespondToGameTypeAction{pass: true, PlayerInfo: pi}
	case 'y':
		return RespondToGameTypeAction{pass: false, PlayerInfo: pi}
	case 'd':
		if len(strMessage) < 5 {
			p.client.send <- []byte("Message Too Short")
			return InvalidAction{PlayerInfo: pi}
		}

		matched, _ := regexp.MatchString("([7-9TQJKA][2345♠♢♡♣]){2,2}", strMessage[1:])

		if !matched {
			p.client.send <- []byte("Message Format Incorrect")
			return InvalidAction{PlayerInfo: pi}
		}

		return ChooseDiscardCardsAction{
			cards: [2]Card{
				stringToCard(strMessage[1:3]),
				stringToCard(strMessage[3:5]),
			},
			PlayerInfo: pi,
		}
	case 's':
		if len(strMessage) != 3 {
			p.client.send <- []byte("Message Length Invalid")
			return InvalidAction{PlayerInfo: pi}
		}

		matched, _ := regexp.MatchString("[7-9TQJKA][2345♠♢♡♣]", strMessage[1:])

		if !matched {
			p.client.send <- []byte("Message Format Incorrect")
			return InvalidAction{PlayerInfo: pi}
		}

		return PlayCardAction{
			card:       stringToCard(strMessage[1:]),
			PlayerInfo: pi,
		}
	default:
		fmt.Println("Invalid action from", p.name)
		return InvalidAction{PlayerInfo: pi}
	}
}


func (p *Player) sendString(message string) {
	p.client.send <- []byte(message)
}