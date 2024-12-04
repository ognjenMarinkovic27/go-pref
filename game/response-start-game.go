package game

type StartGameResponse struct{}

func (r *StartGameResponse) Type() string {
	return "start-game"
}

func (r *StartGameResponse) RecepientPid() string {
	return ""
}
