package network

import "ognjen/go-pref/game"

type ReadyMessage struct {
	MessageBase
}

func (m ReadyMessage) Action() game.Action {
	return game.NewReadyAction(m.Client.player)
}
