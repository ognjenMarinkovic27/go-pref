package network

import "ognjen/go-pref/game"

type InboundMessage interface {
	Action() game.Action
}
