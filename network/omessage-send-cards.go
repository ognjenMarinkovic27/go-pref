package network

import (
	"ognjen/go-pref/game"
)

type SendCardsMesasge struct {
	cards [10]game.Card
	MessageBase
}
