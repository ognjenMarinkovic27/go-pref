package game

type StartHandResponse struct{}

func (r *StartHandResponse) Type() string {
	return "start-hand"
}

func (r *StartHandResponse) RecepientPid() string {
	return ""
}
