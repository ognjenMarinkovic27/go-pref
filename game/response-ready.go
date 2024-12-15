package game

type ReadyResponse struct {
	ReadyPlayerPid string `json:"readyPid"`
}

func (r *ReadyResponse) Type() string {
	return "ready"
}

func (r *ReadyResponse) RecepientPid() string {
	return ""
}
