package game

const (
	Success  = true
	Failiure = false
)

type PlayerResultResponse struct {
	Success   bool
	PlayerPid string
}

func (r *PlayerResultResponse) Type() string {
	return "player-result"
}

func (r *PlayerResultResponse) RecepientPid() string {
	return ""
}
