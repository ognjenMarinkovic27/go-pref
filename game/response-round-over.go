package game

type RoundOverResponse struct {
	Winner *Player `json:"winner"`
}

func (r *RoundOverResponse) Type() string {
	return "round-over"
}

func (r *RoundOverResponse) RecepientPid() string {
	return ""
}
