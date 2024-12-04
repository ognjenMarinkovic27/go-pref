package game

type ReadyAction struct {
	ActionBase `json:"-"`
}

func (action ReadyAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if g.gameState == WaitingGameState && !g.ready[player] {
		return true
	}

	return false
}

func (action ReadyAction) apply(g *Game) {
	player := g.players[action.ppid]
	g.makeReady(player)
	if g.isEveryoneReady() {
		g.startGame(60)
	}
}
