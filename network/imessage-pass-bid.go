package network

import "ognjen/go-pref/game"

type PassBidMessage struct {
	MessageBase
}

func (m PassBidMessage) Action() game.Action {
	return game.NewPassBidAction(m.Client.player)
}
