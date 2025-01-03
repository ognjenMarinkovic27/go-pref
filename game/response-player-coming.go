package game

const (
	Going    = true
	NotGoing = false
)

type PlayerGoingResponse struct {
	Going     bool   `json:"going"`
	PlayerPid string `json:"pid"`
}

func (r *PlayerGoingResponse) Type() string {
	return "player-going"
}

func (r *PlayerGoingResponse) RecepientPid() string {
	return ""
}
