package game

type ReadyAction struct {
	player *Player
}

func NewReadyAction(player *Player) ReadyAction {
	return ReadyAction{player}
}

func (action ReadyAction) validate(g *Game) bool {
	if g.gameState == WaitingGameState && !g.ready[action.player] {
		return true
	}

	return false
}

func (action ReadyAction) apply(g *Game) {
	g.makeReady(action.player)
	if g.isEveryoneReady() {
		g.startGame(60)
	}
}
