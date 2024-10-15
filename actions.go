package main

import (
	"slices"
	"strconv"
)

type Action interface {
	validate(g *Game) bool
	apply(g *Game)
	getPlayer() *Player
}

type PlayerInfo struct {
	player *Player
}

func (pi PlayerInfo) getPlayer() *Player {
	return pi.player
}

type ReadyAction struct {
	PlayerInfo
}

func (action ReadyAction) validate(g *Game) bool {
	if g.gameState == WaitingGameState && !g.ready[action.player] {
		return true
	}

	return false
}

func (action ReadyAction) apply(g *Game) {
	g.makeReady(action.player)
	if g.isEveryoneReady() {
		g.startGame(60)
	}
}

type Bid int

const (
	GameBid       = 0
	TwoBid    Bid = 2
	ThreeBid      = 3
	FourBid       = 4
	FiveBid       = 5
	BattleBid     = 6
	SansBid       = 7
)

type BidAction struct {
	PlayerInfo
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

	g.room.broadcastString("New bid from " + action.player.name + ": " + strconv.Itoa(int(g.currentHandState.bid)))
	g.currentHandState.bidWinner = action.player
	if g.isBiddingWon() {
		g.endBidding()
	} else {
		g.moveToNextActivePlayer()
	}
}

type PassBidAction struct {
	PlayerInfo
}

func (action PassBidAction) validate(g *Game) bool {
	return g.gameState == BiddingGameState && g.isCurrentPlayer(action.player)
}

func (action PassBidAction) apply(g *Game) {
	g.makePlayerPassed(action.player)
	g.room.broadcastString(action.player.name + " passed bidding")

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

type PlayNowBidAction struct {
	PlayerInfo
}

type ChooseGameTypeAction struct {
	gameType GameType
	PlayerInfo
}

func (action ChooseGameTypeAction) validate(g *Game) bool {
	if !g.isCurrentPlayer(action.player) ||
		g.gameState != ChoosingGameTypeGameState ||
		action.gameType < GameType(g.currentHandState.bid) {
		return false
	}

	return true
}

func (action ChooseGameTypeAction) apply(g *Game) {
	g.chooseGameType(action.gameType)
	g.gameState = RespondingToGameTypeGameState
	g.resetPassed()
	g.moveToNextActivePlayer()
}

type RespondToGameTypeAction struct {
	pass bool
	PlayerInfo
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
		g.room.broadcastString(g.currentHandState.currentPlayer.name + "is not coming")
		g.makePlayerPassed(action.player)

		if len(g.currentHandState.passed) == 2 {
			g.room.broadcastString("Nobody is coming! " + g.currentHandState.bidWinner.name + " succeeds!")
			g.currentHandState.bidWinner.score.score -= int(g.currentHandState.gameType) * 2
			g.startNewHand()
			return
		}
	} else {
		g.room.broadcastString(g.currentHandState.currentPlayer.name + " is coming!!!")
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

type ChooseDiscardCardsAction struct {
	cards [2]Card
	PlayerInfo
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

	action.player.sendString("Your new hand:")
	g.sendHandToClient(action.player)
	g.room.broadcastString(action.player.name + " is choosing game type")
	g.transitionToState(ChoosingGameTypeGameState)
}

func otherIndex(index int) int {
	return 1 - index
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

func findCard(card Card, searchSet []Card) int {
	for i := range searchSet {
		if card == searchSet[i] {
			return i
		}
	}

	return -1
}

type PlayCardAction struct {
	card Card
	PlayerInfo
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

	canPlayAppropriateCard := hasAppropriateSuit && g.currentHandState.roundState.suit == action.card.suit
	canPlayTrumpCard := !hasAppropriateSuit && hasTrumpSuit && action.card.suit == trumpSuit
	canPlayAnyCard := !hasAppropriateSuit && !hasTrumpSuit

	if canPlayAppropriateCard || canPlayTrumpCard || canPlayAnyCard {
		return !action.player.played[action.card]
	}

	return false
}

func (p *Player) hasSuit(suit CardSuit) bool {
	for _, card := range p.hand {
		if card.suit == suit && !p.played[card] {
			return true
		}
	}

	return false
}

func (action PlayCardAction) apply(g *Game) {
	g.playCard(action.player, action.card)

	if g.isCurrentRoundOver() {
		p := g.getRoundWinner()
		g.room.broadcastString(p.name + " takes the round")
		g.sendClientsTheirHands()
		g.startNextRound()

		g.currentHandState.roundsPlayed++
		if g.isHandOver() {
			g.room.broadcastString("Hand Done")
			g.checkSuccess()
		}
	} else {
		g.moveToNextActivePlayer()
	}
}

type InvalidAction struct {
	PlayerInfo
}

func (action InvalidAction) validate(g *Game) bool {
	return false
}

func (action InvalidAction) apply(g *Game) {}
