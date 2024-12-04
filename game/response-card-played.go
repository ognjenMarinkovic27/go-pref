package game

type CardPlayedResponse struct {
	Player *Player
	Card   Card
}

func (r *CardPlayedResponse) Type() string {
	return "card-played"
}

func (r *CardPlayedResponse) RecepientPid() string {
	return ""
}
