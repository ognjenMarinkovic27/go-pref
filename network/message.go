package network

type MessageBase struct {
	Seq    int     `json:"seq"`
	Client *Client `json:"-"`
}
