package game

type PlayerPassedResponse struct {
	PasserPid string `json:"pid"`
}

func (r *PlayerPassedResponse) Type() string {
	return "player-passed"
}

func (r *PlayerPassedResponse) RecepientPid() string {
	return ""
}
