package game

type ChoosingCardsResponse struct {
	Chooser     *Player
	HiddenCards [2]Card
}

func (r *ChoosingCardsResponse) Type() string {
	return "choosing-cards"
}

func (r *ChoosingCardsResponse) RecepientPid() string {
	return ""
}
