package game

type RespondToGameTypeAction struct {
	pass   bool
	player *Player
}

func NewRespondToGameTypeAction(pass bool, player *Player) RespondToGameTypeAction {
	return RespondToGameTypeAction{pass, player}
}

func (action RespondToGameTypeAction) validate(g *Game) bool {
	if !g.isCurrentPlayer(action.player) ||
		g.gameState != RespondingToGameTypeGameState {
		return false
	}

	return true
}

func (action RespondToGameTypeAction) apply(g *Game) {
	if action.pass {
		// g.room.broadcastString(g.currentHandState.currentPlayer.getName() + "is not coming")
		g.makePlayerPassed(action.player)

		if len(g.currentHandState.passed) == 2 {
			// g.room.broadcastString("Nobody is coming! " + g.currentHandState.bidWinner.getName() + " succeeds!")
			g.currentHandState.bidWinner.score.score -= int(g.currentHandState.gameType) * 2
			g.startNewHand()
			return
		}
	} else {
		// g.room.broadcastString(g.currentHandState.currentPlayer.getName() + " is coming!!!")
	}

	g.moveToNextActivePlayer()

	if g.isCurrentPlayer(g.currentHandState.bidWinner) {
		if g.currentHandState.bid == SansBid {
			beforePlayer := g.currentHandState.currentPlayer.next.next
			if g.currentHandState.passed[beforePlayer] {
				g.currentHandState.currentPlayer = action.player.next
			} else {
				g.currentHandState.currentPlayer = beforePlayer
			}
		} else {
			g.currentHandState.currentPlayer = g.dealerPlayer.next
		}

		g.transitionToState(PlayingHandGameState)
	}
}
