package game

type StartHandResponse struct {
	FirstPid string `json:"firstPid"`
}

func (r *StartHandResponse) Type() string {
	return "start-hand"
}

func (r *StartHandResponse) RecepientPid() string {
	return ""
}
