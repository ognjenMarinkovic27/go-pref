package network

import (
	"ognjen/go-pref/game"
)

type PlayCardMessage struct {
	Card game.Card
	MessageBase
}

func (m PlayCardMessage) Action() game.Action {
	return game.NewPlayCardAction(m.Card, m.Client.player)
}
