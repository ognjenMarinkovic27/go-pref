package game

type PlayCardAction struct {
	Card       Card `json:"card"`
	ActionBase `json:"-"`
}

func (action PlayCardAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if !g.isCurrentPlayer(player) ||
		g.gameState != PlayingHandGameState {
		return false
	}

	if findCard(action.Card, player.hand[:]) < 0 {
		return false
	}

	if g.currentHandState.roundState.empty {
		return !player.played[action.Card]
	}

	hasAppropriateSuit := player.hasSuit(g.currentHandState.roundState.suit)
	trumpSuit, trumpSuitExists := g.getTrumpSuit()
	hasTrumpSuit := trumpSuitExists && player.hasSuit(trumpSuit)

	canPlayAppropriateCard := hasAppropriateSuit && g.currentHandState.roundState.suit == action.Card.Suit
	canPlayTrumpCard := !hasAppropriateSuit && hasTrumpSuit && action.Card.Suit == trumpSuit
	canPlayAnyCard := !hasAppropriateSuit && !hasTrumpSuit

	if canPlayAppropriateCard || canPlayTrumpCard || canPlayAnyCard {
		return !player.played[action.Card]
	}

	return false
}

func (p *Player) hasSuit(suit CardSuit) bool {
	for _, card := range p.hand {
		if card.Suit == suit && !p.played[card] {
			return true
		}
	}

	return false
}

func (action PlayCardAction) apply(g *Game) {
	player := g.players[action.ppid]
	g.playCard(player, action.Card)

	if g.isCurrentRoundOver() {
		g.reportRoundOver()

		g.currentHandState.roundsPlayed++
		if g.isHandOver() {
			g.reportSuccess()
			g.sendScoresToPlayers()
			g.startNewHand()
		} else {
			g.startNextRound()
		}
	} else {
		g.moveToNextActivePlayer()
	}
}
