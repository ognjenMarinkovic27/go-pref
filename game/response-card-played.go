package game

type CardPlayedResponse struct {
	PlayerPid string `json:"pid"`
	Card      Card   `json:"card"`
}

func (r *CardPlayedResponse) Type() string {
	return "card-played"
}

func (r *CardPlayedResponse) RecepientPid() string {
	return ""
}
