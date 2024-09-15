package main

import (
	"strconv"
)

type Action interface {
	validate(g *Game) bool
	apply(g *Game)
}

type PlayerInfo struct {
	player *Player
}

type ReadyAction struct {
	PlayerInfo
}

func (action ReadyAction) validate(g *Game) bool {
	if g.gameState == WaitingGameState && !g.ready[action.player.id] {
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
		g.currentHandState.firstPlayer = action.player.id
	}

	if !g.isPlayerFirstToBid(action.player) {
		g.currentHandState.currentBid++
	}

	g.room.broadcast <- []byte("New bid from " + action.player.name + ": " + strconv.Itoa(int(g.currentHandState.currentBid)))
	g.moveToNextActivePlayer()
}

func (action PassBidAction) validate(g *Game) bool {
	return g.gameState == BiddingGameState && g.isCurrentPlayer(action.player)
}

func (action PassBidAction) apply(g *Game) {
	g.makePlayerPassed(action.player)
	if g.isBiddingWon() {
		g.transitionToChooseGameType()
	}

	if g.isEveryonePassed() {
		g.startNewHand()
	}

	g.moveToNextActivePlayer()
}

type BidAction struct {
	PlayerInfo
}

type PassBidAction struct {
	PlayerInfo
}

type PlayNowBidAction struct {
	PlayerInfo
}

type PlayCardAction struct {
	card Card
	PlayerInfo
}

type ChooseDiscardCardsAction struct {
	card1 Card
	card2 Card
	PlayerInfo
}

type ChooseGameTypeAction struct {
	gameType GameType
	PlayerInfo
}
