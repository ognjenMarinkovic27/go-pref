package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
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
	firstPlayer       *Player
	currentPlayer     *Player
	bidWinner         *Player
	currentBid        Bid
	passed            map[*Player]bool
	currentGameType   GameType
	currentRoundState RoundState
	hiddenCards       [2]Card
}

type RoundState struct {
	empty       bool
	currentSuit CardSuit
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
		slices.SortFunc(p.hand[:], func(a, b Card) int {
			if a.suit-b.suit == 0 {
				return int(b.value - a.value)
			} else {
				return int(a.suit - b.suit)
			}
		})
		s += 10
	}

	g.currentHandState.hiddenCards[0] = deck[30]
	g.currentHandState.hiddenCards[1] = deck[31]
}

func (g *Game) run() {
	for {
		fmt.Println("GAme")
		if g.started {
			g.room.broadcast <- []byte("Turn: " + g.currentHandState.currentPlayer.name)
		}
		action := <-g.actions
		if g.validate(action) {
			g.apply(action)
		} else {
			action.getPlayer().client.send <- []byte("Invalid action")
		}
	}
}

func (g *Game) makeReady(p *Player) {
	g.ready[p] = true
	g.room.broadcast <- []byte(p.name + "is ready!")
}

func (g *Game) isEveryoneReady() bool {
	return len(g.ready) == 3
}

func (g *Game) startGame() {
	g.room.broadcast <- []byte("Starting game")
	g.started = true
	g.setupPlayers()
	g.startNewHand()
}

func (g *Game) setupPlayers() {
	var prev *Player = nil
	var first *Player = nil
	count := 0
	for p := range g.players {
		if prev != nil {
			prev.next = p
		} else {
			g.dealerPlayer = p
			first = p
		}

		count++
		if count == 3 {
			p.next = first
		}

		prev = p
	}
}

func (g *Game) startNewHand() {
	g.room.broadcast <- []byte("Starting hand")

	g.gameState = BiddingGameState
	g.currentHandState = HandState{
		firstPlayer:     nil,
		currentPlayer:   g.nextPlayer(g.dealerPlayer),
		passed:          make(map[*Player]bool),
		currentBid:      TwoBid,
		currentGameType: NoneGameType,
	}

	g.dealerPlayer = g.nextPlayer(g.dealerPlayer)

	g.dealCards()
	g.sendClientsTheirHands()
}

func (g *Game) sendClientsTheirHands() {
	for p := range g.players {
		str := ""

		for i, c := range p.hand {
			str += cardToString(c)
			if i != 10 {
				str += " "
			}
		}

		p.client.send <- []byte(str)
	}
}

func (g *Game) isCurrentPlayer(p *Player) bool {
	return g.currentHandState.currentPlayer == p
}

func (g *Game) isBiddingMaxed() bool {
	return g.currentHandState.currentBid == SansBid
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
	g.room.broadcast <- []byte(p.name + " passed")
}

func (g *Game) isBiddingWon() bool {
	return len(g.currentHandState.passed) == 2 && g.currentHandState.firstPlayer != nil
}

func (g *Game) transitionToState(state GameState) {
	g.gameState = state
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
