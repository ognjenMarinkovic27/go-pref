package game

type PassBidAction struct {
	ActionBase `json:"-"`
}

func (action PassBidAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	return g.gameState == BiddingGameState && g.isCurrentPlayer(player)
}

func (action PassBidAction) apply(g *Game) {
	player := g.players[action.ppid]
	g.makePlayerPassed(player)

	if g.isBiddingWon() {
		g.endBidding()
		return
	}

	if g.isEveryonePassed() {
		g.startNewHand()
		return
	}

	g.moveToNextActivePlayer()
}
