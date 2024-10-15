package main

import (
	"math/rand/v2"
	"slices"
	"strconv"
)

type GameState int

const (
	WaitingGameState              GameState = 0
	BiddingGameState              GameState = 1
	PlayNowGameState              GameState = 2
	ClaimPlayNowTypeGameState     GameState = 3
	ChoosingGameTypeGameState     GameState = 4
	RespondingToGameTypeGameState GameState = 5
	ChoosingCardsGameState        GameState = 6
	PlayingHandGameState          GameState = 7
)

type GameType int

const (
	NoneGameType     GameType = 0
	SpadesGameType   GameType = 2
	DiamondsGameType GameType = 3
	HeartsGameType   GameType = 4
	ClubsGameType    GameType = 5
	BattleGameType   GameType = 6
	SansGameType     GameType = 7
)

var gameTypeToSuit = map[GameType]CardSuit{
	SpadesGameType:   Spades,
	DiamondsGameType: Diamonds,
	HeartsGameType:   Hearts,
	ClubsGameType:    Clubs,
}

type Game struct {
	gameState        GameState
	dealerPlayer     *Player
	currentHandState HandState
	players          map[*Player]bool
	actions          chan Action

	ready   map[*Player]bool
	started bool

	room *Room
}

type HandState struct {
	firstPlayer   *Player
	currentPlayer *Player
	bidWinner     *Player
	bid           Bid
	passed        map[*Player]bool
	gameType      GameType
	roundsPlayed  int
	roundState    RoundState
	roundsWon     map[*Player]int
	hiddenCards   [2]Card
}

type RoundState struct {
	empty bool
	suit  CardSuit
	table map[*Player]Card
}

func newGame(actions chan Action, room *Room) *Game {
	return &Game{
		gameState: WaitingGameState,
		players:   make(map[*Player]bool),
		actions:   actions,
		ready:     make(map[*Player]bool),
		started:   false,
		room:      room,
	}
}

func (g *Game) addPlayer(p *Player) bool {
	g.players[p] = true
	return true
}

func (g *Game) removePlayer(p *Player) {
	if _, ok := g.players[p]; ok {
		delete(g.players, p)
	}
}

func (g *Game) validate(a Action) bool {
	if a == nil {
		return false
	}

	return a.validate(g)
}

func (g *Game) apply(a Action) {
	a.apply(g)
}

func (g *Game) dealCards() {
	var deck [32]Card

	c := 0
	for i := Seven; i <= Ace; i++ {
		for j := Spades; j <= Hearts; j++ {
			deck[c] = Card{suit: j, value: i}
			c++
		}
	}

	for i := 31; i > 0; i-- {
		j := rand.IntN(i)

		deck[i], deck[j] = deck[j], deck[i]
	}

	s := 0
	for p := range g.players {
		copy(p.hand[:], deck[s:s+10])
		slices.SortFunc(p.hand[:], cardCompare)
		s += 10
	}

	g.currentHandState.hiddenCards[0] = deck[30]
	g.currentHandState.hiddenCards[1] = deck[31]
}

func (g *Game) run() {
	for {
		if g.started {
			g.room.broadcastString("Turn: " + g.getCurrentPlayer().name)
		}
		action := <-g.actions
		g.handleAction(action)
	}
}

func (g *Game) handleAction(action Action) {
	if g.validate(action) {
		g.apply(action)
	} else {
		action.getPlayer().client.send <- []byte("Invalid action")
	}
}

func (g *Game) getCurrentPlayer() *Player {
	return g.currentHandState.currentPlayer
}

func (g *Game) sendClientsTheirScores() {
	for p := range g.players {
		p.sendString("Score " + strconv.Itoa(p.score.score))

		for other, soup := range p.score.soups {
			p.sendString(other.name + ": " + strconv.Itoa(soup))
		}
	}
}

func (g *Game) makeReady(p *Player) {
	g.ready[p] = true
	g.room.broadcastString(p.name + "is ready!")
}

func (g *Game) isEveryoneReady() bool {
	return len(g.ready) == 3
}

func (g *Game) startGame(startingScore int) {
	g.room.broadcastString("Starting game")
	g.started = true
	g.setupPlayers(startingScore)
	g.startNewHand()
}

func (g *Game) setupPlayers(startingScore int) {
	var prev *Player = nil
	var first *Player = nil
	count := 0
	for p := range g.players {
		if prev != nil {
			prev.next = p
			p.score.soups[prev] = 0
			prev.score.soups[p] = 0
		} else {
			g.dealerPlayer = p
			first = p
		}

		p.score.score = startingScore

		count++
		if count == 3 {
			p.next = first
			p.score.soups[first] = 0
			first.score.soups[p] = 0
		}

		prev = p
	}
}

func (g *Game) startNewHand() {
	g.room.broadcastString("Starting hand")

	g.gameState = BiddingGameState
	g.currentHandState = HandState{
		firstPlayer:   nil,
		currentPlayer: g.nextPlayer(g.dealerPlayer),
		passed:        make(map[*Player]bool),
		bid:           TwoBid,
		gameType:      NoneGameType,
		roundsPlayed:  0,
		roundState: RoundState{
			empty: true,
			table: make(map[*Player]Card),
		},
		roundsWon: make(map[*Player]int),
	}

	g.dealerPlayer = g.nextPlayer(g.dealerPlayer)

	g.dealCards()
	g.sendClientsTheirHands()
	g.sendClientsTheirScores()
}

func (g *Game) sendClientsTheirHands() {
	for p := range g.players {
		g.sendHandToClient(p)
	}
}

func (g *Game) sendHandToClient(p *Player) {
	str := ""

	for i, c := range p.hand {
		if p.played[c] {
			continue
		}

		str += cardToString(c)
		if i != 10 {
			str += " "
		}
	}

	p.client.send <- []byte(str)
}

func (g *Game) isCurrentPlayer(p *Player) bool {
	return g.currentHandState.currentPlayer == p
}

func (g *Game) isBiddingMaxed() bool {
	return g.currentHandState.bid == SansBid
}

func (g *Game) hasPassed(p *Player) bool {
	return g.currentHandState.passed[p]
}

func (g *Game) isFirstBid() bool {
	return g.currentHandState.firstPlayer == nil
}

func (g *Game) isPlayerFirstToBid(p *Player) bool {
	return g.currentHandState.firstPlayer == p
}

func (g *Game) chooseGameType(gameType GameType) {
	g.currentHandState.gameType = gameType
	trumpSuit, exists := gameTypeToSuit[gameType]
	if exists {
		g.currentHandState.roundState.suit = trumpSuit
	}

	g.room.broadcastString(g.currentHandState.currentPlayer.name + " chose game type: " + strconv.Itoa(int(gameType)))
}

func (g *Game) resetPassed() {
	clear(g.currentHandState.passed)
}

func (g *Game) moveToNextActivePlayer() {
	g.currentHandState.currentPlayer = g.nextPlayer(g.currentHandState.currentPlayer)

	if g.currentHandState.passed[g.currentHandState.currentPlayer] {
		g.currentHandState.currentPlayer = g.nextPlayer(g.currentHandState.currentPlayer)
	}
}

func (g *Game) makePlayerPassed(p *Player) {
	g.currentHandState.passed[p] = true
	g.room.broadcastString(p.name + " passed")
}

func (g *Game) isBiddingWon() bool {
	return len(g.currentHandState.passed) == 2 && g.currentHandState.firstPlayer != nil
}

func (g *Game) transitionToState(state GameState) {
	g.gameState = state
}

func (g *Game) makeNonPassedPlayerCurrent() {
	for p := range g.players {
		if !g.currentHandState.passed[p] {
			g.currentHandState.currentPlayer = p
			break
		}
	}
}

func (g *Game) isEveryonePassed() bool {
	return len(g.currentHandState.passed) == 3
}

func (g *Game) nextPlayer(p *Player) *Player {
	return p.next
}

func (g *Game) playCard(p *Player, card Card) {
	g.currentHandState.roundState.table[p] = card
	if (g.currentHandState.roundState.empty) {
		g.currentHandState.roundState.suit = card.suit
	}
	g.currentHandState.roundState.empty = false
	
	p.played[card] = true

	g.room.broadcastString(p.name + " played " + cardToString(card))
}

func (g *Game) isCurrentRoundOver() bool {
	return len(g.currentHandState.roundState.table) == 3-len(g.currentHandState.passed)
}

func (g *Game) getTrumpSuit() (CardSuit, bool) {
	cs, ok := gameTypeToSuit[g.currentHandState.gameType]
	return cs, ok
}

func (g *Game) getRoundWinner() *Player {
	var roundWinner *Player = nil
	var bestCard Card
	for p, c := range g.currentHandState.roundState.table {
		if roundWinner == nil {
			roundWinner = p
			bestCard = c
		} else {
			trump, _ := g.getTrumpSuit()
			if (c.suit == trump && (bestCard.suit != trump || bestCard.value < c.value)) ||
				c.suit != trump && bestCard.value < c.value {
				roundWinner = p
				bestCard = c
			}
		}
	}

	return roundWinner
}

func (g *Game) startNextRound() {
	p := g.getRoundWinner()
	g.currentHandState.roundState.empty = true
	clear(g.currentHandState.roundState.table)
	g.currentHandState.currentPlayer = p
	g.currentHandState.roundsWon[p]++
}

func (g *Game) isHandOver() bool {
	return g.currentHandState.roundsPlayed == 10
}

func (g *Game) checkSuccess() {
	owner := g.currentHandState.bidWinner
	if g.currentHandState.roundsWon[owner] >= 6 {
		owner.score.score -= int(g.currentHandState.gameType) * 2
		g.room.broadcastString(owner.name + " succeded")
	} else {
		owner.score.score += int(g.currentHandState.gameType) * 2
		g.room.broadcastString(owner.name + " failed")
	}

	for p := range g.players {
		if p == owner {
			continue
		}

		if g.currentHandState.roundsWon[p] >= 2 || 10-g.currentHandState.roundsWon[owner] <= 6 {
			g.room.broadcastString(p.name + " succeded")
		} else {
			owner.score.score += int(g.currentHandState.gameType) * 2
			g.room.broadcastString(p.name + " failed :(")
		}

		p.score.soups[owner] += g.currentHandState.roundsWon[p] * int(g.currentHandState.gameType) * 2
	}
}
