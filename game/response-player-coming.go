package game

const (
	Coming    = true
	NotComing = false
)

type PlayerGoingResponse struct {
	Coming    bool
	PlayerPid string
}

func (r *PlayerGoingResponse) Type() string {
	return "player-going"
}

func (r *PlayerGoingResponse) RecepientPid() string {
	return ""
}
