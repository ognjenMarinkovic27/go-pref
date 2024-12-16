package game

type BidAction struct {
	ActionBase `json:"-"`
}

func (action BidAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if g.gameState == BiddingGameState &&
		g.isCurrentPlayer(player) &&
		!g.isBiddingMaxed() &&
		!g.hasPassed(player) {
		return true
	}

	return false
}

func (action BidAction) apply(g *Game) {
	player := g.players[action.ppid]
	if g.isFirstBid() {
		g.currentHandState.firstBidder = player
	}

	if !g.isPlayerFirstToBid(player) {
		g.currentHandState.bid++
	}

	g.reportBid(player)
	g.currentHandState.bidWinner = player
	if g.isBiddingWon() {
		g.endBidding()
	} else {
		g.moveToNextActivePlayer()
	}
}
