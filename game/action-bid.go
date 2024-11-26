package game

type BidAction struct {
	player *Player
}

func NewBidAction(player *Player) BidAction {
	return BidAction{player}
}

func (action BidAction) validate(g *Game) bool {
	if g.gameState == BiddingGameState &&
		g.isCurrentPlayer(action.player) &&
		!g.isBiddingMaxed() &&
		!g.hasPassed(action.player) {
		return true
	}

	return false
}

func (action BidAction) apply(g *Game) {
	if g.isFirstBid() {
		g.currentHandState.firstPlayer = action.player
	}

	if !g.isPlayerFirstToBid(action.player) {
		g.currentHandState.bid++
	}

	// g.room.broadcastString("New bid from " + action.player.getName() + ": " + strconv.Itoa(int(g.currentHandState.bid)))
	g.currentHandState.bidWinner = action.player
	if g.isBiddingWon() {
		g.endBidding()
	} else {
		g.moveToNextActivePlayer()
	}
}
