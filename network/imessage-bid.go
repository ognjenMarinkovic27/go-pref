package network

import (
	"ognjen/go-pref/game"
)

type BidMessage struct {
	MessageBase
}

func (m BidMessage) Action() game.Action {
	return game.NewBidAction(m.Client.player)
}
