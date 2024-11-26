package game

type PassBidAction struct {
	player *Player
}

func NewPassBidAction(player *Player) PassBidAction {
	return PassBidAction{player}
}

func (action PassBidAction) validate(g *Game) bool {
	return g.gameState == BiddingGameState && g.isCurrentPlayer(action.player)
}

func (action PassBidAction) apply(g *Game) {
	g.makePlayerPassed(action.player)

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
