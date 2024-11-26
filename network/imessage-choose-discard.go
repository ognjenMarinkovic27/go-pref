package network

import (
	"ognjen/go-pref/game"
)

type ChooseDiscardCardsMessage struct {
	Cards [2]game.Card `json:"cards"`
	MessageBase
}

func (m ChooseDiscardCardsMessage) Action() game.Action {
	return game.NewChooseDiscardCardsAction(m.Cards, m.Client.player)
}
