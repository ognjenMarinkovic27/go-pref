package game

import "slices"

type ChooseDiscardCardsAction struct {
	Cards      [2]Card `json:"cards"`
	ActionBase `json:"-"`
}

func (action ChooseDiscardCardsAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if !g.isCurrentPlayer(player) ||
		g.gameState != ChoosingCardsGameState {
		return false
	}

	var found [2]bool

	foundInHand := containsCards(action.Cards[:], player.hand[:])
	foundInHidden := containsCards(action.Cards[:], g.currentHandState.hiddenCards[:])

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
	player := g.players[action.ppid]
	for _, c := range action.Cards {
		index := findCard(c, player.hand[:])
		if index < 0 {
			continue
		}

		swapIndex := findDifferentThan(action.Cards[:], g.currentHandState.hiddenCards[:])

		swapCards(&player.hand[index], &g.currentHandState.hiddenCards[swapIndex])
	}

	slices.SortFunc(player.hand[:], cardCompare)

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
