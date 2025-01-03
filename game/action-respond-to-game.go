package game

type RespondToGameTypeAction struct {
	Pass bool `json:"pass"`
	ActionBase
}

func (action RespondToGameTypeAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if !g.isCurrentPlayer(player) ||
		g.gameState != RespondingToGameTypeGameState {
		return false
	}

	return true
}

func (action RespondToGameTypeAction) apply(g *Game) {
	player := g.players[action.ppid]
	if action.Pass {
		g.recordPlayerGoingState(player, NotGoing)
		g.currentHandState.passed[player] = true

		if len(g.currentHandState.passed) == 2 {
			g.currentHandState.bidWinner.score.main -= int(g.currentHandState.gameType) * 2
			g.reportSuccessToOwner()
			g.startNewHand()
			return
		}
	} else {
		g.recordPlayerGoingState(player, Going)
	}

	g.moveToNextActivePlayer()

	if g.isCurrentPlayer(g.currentHandState.bidWinner) {
		if g.currentHandState.bid == SansBid {
			beforePlayer := g.currentHandState.currentPlayer.next.next
			if g.currentHandState.passed[beforePlayer] {
				g.currentHandState.currentPlayer = player.next
			} else {
				g.currentHandState.currentPlayer = beforePlayer
			}
		} else {
			g.currentHandState.currentPlayer = g.dealerPlayer
		}

		g.transitionToState(PlayingHandGameState)
	}
}
