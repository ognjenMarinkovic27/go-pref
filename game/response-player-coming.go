package game

const (
	Coming    = true
	NotComing = false
)

type PlayerComingResponse struct {
	Coming    bool
	PlayerPid string
}

func (r *PlayerComingResponse) Type() string {
	return "player-coming"
}

func (r *PlayerComingResponse) RecepientPid() string {
	return ""
}
