package game

import "slices"

type ChooseDiscardCardsAction struct {
	cards  [2]Card
	player *Player
}

func NewChooseDiscardCardsAction(cards [2]Card, player *Player) ChooseDiscardCardsAction {
	return ChooseDiscardCardsAction{cards, player}
}

func (action ChooseDiscardCardsAction) validate(g *Game) bool {
	if !g.isCurrentPlayer(action.player) ||
		g.gameState != ChoosingCardsGameState {
		return false
	}

	var found [2]bool

	foundInHand := containsCards(action.cards[:], action.player.hand[:])
	foundInHidden := containsCards(action.cards[:], g.currentHandState.hiddenCards[:])

	for i := range found {
		found[i] = foundInHand[i] || foundInHidden[i]
	}

	return found[0] && found[1]
}

func containsCards(cards []Card, searchSet []Card) (found [2]bool) {
	for _, card := range searchSet {
		index := findCard(card, cards)

		if index >= 0 {
			found[index] = true
		}
	}

	return
}

func (action ChooseDiscardCardsAction) apply(g *Game) {
	for _, c := range action.cards {
		index := findCard(c, action.player.hand[:])
		if index < 0 {
			continue
		}

		swapIndex := findDifferentThan(action.cards[:], g.currentHandState.hiddenCards[:])

		swapCards(&action.player.hand[index], &g.currentHandState.hiddenCards[swapIndex])
	}

	slices.SortFunc(action.player.hand[:], cardCompare)

	g.transitionToState(ChoosingGameTypeGameState)
}

func findDifferentThan(dontMatchCards []Card, inSet []Card) int {
in_set_iter:
	for ind, card := range inSet {
		for _, dontMatchCard := range dontMatchCards {
			if card == dontMatchCard {
				continue in_set_iter
			}
		}

		return ind
	}

	return -1
}

func swapCards(card1, card2 *Card) {
	*card1, *card2 = *card2, *card1
}
