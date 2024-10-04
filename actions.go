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
		g.startGame()
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
		g.currentHandState.currentBid++
	}

	g.room.broadcast <- []byte("New bid from " + action.player.name + ": " + strconv.Itoa(int(g.currentHandState.currentBid)))
	g.currentHandState.bidWinner = action.player
	g.moveToNextActivePlayer()
}

type PassBidAction struct {
	PlayerInfo
}

func (action PassBidAction) validate(g *Game) bool {
	return g.gameState == BiddingGameState && g.isCurrentPlayer(action.player)
}

func (action PassBidAction) apply(g *Game) {
	g.makePlayerPassed(action.player)
	g.room.broadcast <- []byte(action.player.name + " passed bidding")

	if g.isBiddingWon() {
		g.transitionToState(ChoosingGameTypeGameState)
		g.makeNonPassedPlayerCurrent()
		g.room.broadcast <- []byte(g.currentHandState.bidWinner.name + " is choosing game type")
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
		action.gameType < GameType(g.currentHandState.currentBid) {
		return false
	}

	return true
}

func (action ChooseGameTypeAction) apply(g *Game) {
	g.currentHandState.currentGameType = action.gameType
	g.room.broadcast <- []byte(g.currentHandState.currentPlayer.name + " chose game type: " + strconv.Itoa(int(action.gameType)))
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
		g.room.broadcast <- []byte(g.currentHandState.currentPlayer.name + "is not coming")
		g.makePlayerPassed(action.player)
	} else {
		g.room.broadcast <- []byte(g.currentHandState.currentPlayer.name + "is coming!!!")
	}

	g.moveToNextActivePlayer()

	if g.isCurrentPlayer(g.currentHandState.bidWinner) {
		g.transitionToState(ChoosingCardsGameState)
		g.room.broadcast <- []byte(g.currentHandState.bidWinner.name + g.currentHandState.currentPlayer.name + " is choosing cards")
		g.currentHandState.bidWinner.client.send <- []byte("Hidden cards: " +
			cardToString(g.currentHandState.hiddenCards[0]) + " " +
			cardToString(g.currentHandState.hiddenCards[1]))
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

	foundInHand := containsCards(action.cards[:], g.currentHandState.currentPlayer.hand[:])
	foundInHidden := containsCards(action.cards[:], g.currentHandState.hiddenCards[:])

	for i := range found {
		found[i] = foundInHand[i] || foundInHidden[i]
	}

	println(foundInHand[0], foundInHand[1], foundInHidden[0], foundInHidden[1])
	println(cardToString(action.cards[0]), cardToString(action.cards[1]))

	return found[0] && found[1]
}

func containsCards(cards []Card, searchSet []Card) (found [2]bool) {
	for _, card := range searchSet {
		index := findCard(card, cards)

		if index > 0 {
			found[index] = true
		}
	}

	return
}

func (action ChooseDiscardCardsAction) apply(g *Game) {
	swapIndex := 0
	for i := range g.currentHandState.hiddenCards {
		if findCard(g.currentHandState.hiddenCards[i], action.cards[:]) > 0 {
			swapIndex++
		}
	}

	if swapIndex >= len(g.currentHandState.hiddenCards) {
		return
	}

	for i := range action.cards {
		if swapIndex >= len(g.currentHandState.hiddenCards) {
			break
		}

		ci := findCard(action.cards[i], action.player.hand[:])
		 
		if ci > 0 {
			action.player.hand[ci], action.cards[i] = 
				action.cards[i], action.player.hand[ci]
			swapIndex++
		}
	}

	slices.SortFunc(action.player.hand[:], func(a, b Card) int {
		if a.suit-b.suit == 0 {
			return int(b.value - a.value)
		} else {
			return int(a.suit - b.suit)
		}
	})

	action.player.client.send <- []byte("Your new hand:")
	g.sendHandToClient(action.player)
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

	return true
}

type InvalidAction struct {
	PlayerInfo
}

func (action InvalidAction) validate(g *Game) bool {
	return false
}

func (action InvalidAction) apply(g *Game) {}
