package game

type InvalidAction struct {
	player *Player
}

func NewInvalidAction(player *Player) InvalidAction {
	return InvalidAction{player}
}

func (action InvalidAction) validate(g *Game) bool {
	return false
}

func (action InvalidAction) apply(g *Game) {}
