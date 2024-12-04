package network

type ConnectedPayload struct {
	Pid string `json:"pid"`
}

func (p *ConnectedPayload) Type() string {
	return "player-connected"
}

type DisconnectedPayload struct {
	Pid string `json:"pid"`
}

func (p *DisconnectedPayload) Type() string {
	return "player-disconnected"
}

type LobbyPayload struct {
	Pids []string `json:"pids"`
}

func (p *LobbyPayload) Type() string {
	return "lobby-state"
}
