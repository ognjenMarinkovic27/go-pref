package network

import (
	"ognjen/go-pref/game"
)

type ChooseGameTypeMessage struct {
	GameType game.GameType `json:"game-type"`
	MessageBase
}

func (m ChooseGameTypeMessage) Action() game.Action {
	return game.NewChooseGameTypeAction(m.GameType, m.Client.player)
}
