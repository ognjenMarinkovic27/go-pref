package network

import (
	"encoding/json"
)

type Payload interface {
	Type() string
}

type OutboundMessage struct {
	Recepient *Client
	Seq       int
	Payload   Payload
}

type rawMessage struct {
	Recepient *Client     `json:"-"`
	Seq       int         `json:"seq"`
	Payload   interface{} `json:"payload"`
}

func (m OutboundMessage) MarshalJSON() ([]byte, error) {
	typedPayload := struct {
		PayloadType string  `json:"payload-type"`
		Payload     Payload `json:"payload"`
	}{
		PayloadType: m.Payload.Type(),
		Payload:     m.Payload,
	}

	rmsg := rawMessage{
		Recepient: m.Recepient,
		Seq:       m.Seq,
		Payload:   typedPayload,
	}

	return json.Marshal(rmsg)
}
