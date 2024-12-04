package network

type InboundMessage struct {
	Seq     int         `json:"seq"`
	Payload interface{} `json:"payload"`
	Client  *Client     `json:"-"`
}
