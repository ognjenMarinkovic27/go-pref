package game

type ReadyResponse struct {
	ReadyPlayerPid string
}

func (r *ReadyResponse) Type() string {
	return "ready"
}

func (r *ReadyResponse) RecepientPid() string {
	return r.ReadyPlayerPid
}
