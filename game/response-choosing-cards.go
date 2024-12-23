package game

type ChoosingCardsResponse struct {
	ChooserPid  string  `json:"chooserPid"`
	HiddenCards [2]Card `json:"hiddenCards"`
}

func (r *ChoosingCardsResponse) Type() string {
	return "choosing-cards"
}

func (r *ChoosingCardsResponse) RecepientPid() string {
	return ""
}
