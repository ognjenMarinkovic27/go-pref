package network

import "ognjen/go-pref/game"

type ResponseAdapter interface {
	toOutboundMessage(response game.Response) OutboundMessage
}
