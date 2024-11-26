package game

type ChoosingCardsResponse struct {
	Chooser     *Player
	HiddenCards [2]Card
}
