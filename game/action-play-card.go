package game

type PlayCardAction struct {
	card   Card
	player *Player
}

func NewPlayCardAction(card Card, player *Player) PlayCardAction {
	return PlayCardAction{card, player}
}

func (action PlayCardAction) validate(g *Game) bool {
	if !g.isCurrentPlayer(action.player) ||
		g.gameState != PlayingHandGameState {
		return false
	}

	if findCard(action.card, action.player.hand[:]) < 0 {
		return false
	}

	if g.currentHandState.roundState.empty {
		return !action.player.played[action.card]
	}

	hasAppropriateSuit := action.player.hasSuit(g.currentHandState.roundState.suit)
	trumpSuit, trumpSuitExists := g.getTrumpSuit()
	hasTrumpSuit := trumpSuitExists && action.player.hasSuit(trumpSuit)

	canPlayAppropriateCard := hasAppropriateSuit && g.currentHandState.roundState.suit == action.card.Suit
	canPlayTrumpCard := !hasAppropriateSuit && hasTrumpSuit && action.card.Suit == trumpSuit
	canPlayAnyCard := !hasAppropriateSuit && !hasTrumpSuit

	if canPlayAppropriateCard || canPlayTrumpCard || canPlayAnyCard {
		return !action.player.played[action.card]
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
	g.playCard(action.player, action.card)

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
