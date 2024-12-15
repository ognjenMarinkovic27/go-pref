package game

type StartGameResponse struct {
	PidOrder []string `json:"pidOrder"`
}

func (r *StartGameResponse) Type() string {
	return "start-game"
}

func (r *StartGameResponse) RecepientPid() string {
	return ""
}
