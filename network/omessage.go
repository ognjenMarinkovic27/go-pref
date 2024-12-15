package network

import (
	"encoding/json"
)

type Payload interface {
	Type() string
}

type OutboundMessage struct {
	Recepient *Client `json:"-"`
	Seq       int     `json:"seq"`
	Payload   Payload `json:"payload"`
}

func (m OutboundMessage) MarshalJSON() ([]byte, error) {
	type Alias OutboundMessage

	type TypedOutboundMessage struct {
		Type string `json:"type"`
		Alias
	}

	typedMessage := TypedOutboundMessage{
		Type:  m.Payload.Type(),
		Alias: (Alias)(m),
	}

	return json.Marshal(typedMessage)
}
