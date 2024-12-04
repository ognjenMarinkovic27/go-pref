package game

type PlayerPassedResponse struct {
	PasserPid string
}

func (r *PlayerPassedResponse) Type() string {
	return "player-passed"
}

func (r *PlayerPassedResponse) RecepientPid() string {
	return ""
}
