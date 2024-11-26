package game

type InvalidAction struct {
	player *Player
}

func NewInvalidAction(player *Player) InvalidAction {
	return InvalidAction{player}
}

func (action InvalidAction) validate(g *Game) bool {
	// We will allow the InvalidAction to pass so we
	// can add an InvalidActionResponse in apply
	return true
}

func (action InvalidAction) apply(g *Game) {
	g.addResponse(&InvalidActionResponse{action.player})
}
