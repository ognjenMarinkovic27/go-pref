package game

import (
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

	responses []Response

	ready   map[*Player]bool
	started bool
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

func NewGame() *Game {
	return &Game{
		gameState: WaitingGameState,
		players:   make(map[*Player]bool),
		ready:     make(map[*Player]bool),
		started:   false,
	}
}

func (g *Game) AddPlayer(p *Player) bool {
	g.players[p] = true
	return true
}

func (g *Game) RemovePlayer(p *Player) {
	delete(g.players, p)
}

func (g *Game) Validate(a Action) bool {
	if a == nil {
		return false
	}

	return a.validate(g)
}

func (g *Game) Apply(a Action) {
	a.apply(g)
}

func (g *Game) Collect() []Response {
	defer func() {
		g.responses = nil
	}()

	return g.responses
}

func (g *Game) dealCards() {
	var deck [32]Card

	c := 0
	for i := Seven; i <= Ace; i++ {
		for j := Spades; j <= Hearts; j++ {
			deck[c] = Card{Suit: j, Value: i}
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

func (g *Game) currentPlayer() *Player {
	return g.currentHandState.currentPlayer
}

func (g *Game) sendScoresToPlayers() {
	for p := range g.players {
		g.addResponse(&SendScoresResponse{p})
	}
}

func (g *Game) makeReady(p *Player) {
	g.ready[p] = true
	g.addResponse(&ReadyResponse{p})
}

func (g *Game) isEveryoneReady() bool {
	return len(g.ready) == 3
}

func (g *Game) startGame(startingScore int) {
	g.addResponse(&StartGameResponse{})
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

		p.score.main = startingScore

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
	g.addResponse(&StartHandResponse{})
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

	g.clearPlayedMaps()
	g.dealCards()

	g.sendHandsToPlayers()
	g.sendScoresToPlayers()
}

func (g *Game) sendHandsToPlayers() {
	for p := range g.players {
		g.sendHandToPlayer(p)
	}
}

func (g *Game) sendHandToPlayer(p *Player) {
	g.addResponse(&SendCardsResponse{p})
}

func (g *Game) clearPlayedMaps() {
	for p := range g.players {
		clear(p.played)
	}
}

func (g *Game) isCurrentPlayer(p *Player) bool {
	return g.currentHandState.currentPlayer == p
}

func (g *Game) reportBid(p *Player) {
	g.addResponse(&NewBidResponse{p, g.currentHandState.bid})
}

func (g *Game) isBiddingMaxed() bool {
	return g.currentHandState.bid == SansBid
}

func (g *Game) endBidding() {
	g.transitionToState(ChoosingCardsGameState)
	g.makeNonPassedPlayerCurrent()
	g.addResponse(&ChoosingCardsResponse{g.currentPlayer(), g.currentHandState.hiddenCards})
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

	g.addResponse(&GameTypeChosenResponse{gameType})
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

func (g *Game) recordPlayerComingState(p *Player, coming bool) {
	g.addResponse(&PlayerComingResponse{coming, p})
}

func (g *Game) makePlayerPassed(p *Player) {
	g.currentHandState.passed[p] = true
	g.addResponse(&PlayerPassedResponse{p})
}

func (g *Game) isBiddingWon() bool {
	return len(g.currentHandState.passed) == 2 && g.currentHandState.firstPlayer != nil
}

func (g *Game) transitionToState(state GameState) {
	g.gameState = state
	g.addResponse(&GameStateResponse{g.gameState})
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
	if g.currentHandState.roundState.empty {
		g.currentHandState.roundState.suit = card.Suit
	}
	g.currentHandState.roundState.empty = false

	p.played[card] = true

	g.addResponse(&CardPlayedResponse{p, card})
}

func (g *Game) isCurrentRoundOver() bool {
	return len(g.currentHandState.roundState.table) == 3-len(g.currentHandState.passed)
}

func (g *Game) reportRoundOver() {
	g.addResponse(&RoundOverResponse{g.roundWinner()})
}

func (g *Game) getTrumpSuit() (CardSuit, bool) {
	cs, ok := gameTypeToSuit[g.currentHandState.gameType]
	return cs, ok
}

func (g *Game) roundWinner() *Player {
	var roundWinner *Player = nil
	var bestCard Card
	for p, c := range g.currentHandState.roundState.table {
		if roundWinner == nil {
			roundWinner = p
			bestCard = c
		} else {
			trump, _ := g.getTrumpSuit()
			if (c.Suit == trump && (bestCard.Suit != trump || bestCard.Value < c.Value)) ||
				c.Suit != trump && bestCard.Value < c.Value {
				roundWinner = p
				bestCard = c
			}
		}
	}

	return roundWinner
}

func (g *Game) startNextRound() {
	p := g.roundWinner()
	g.currentHandState.roundState.empty = true
	clear(g.currentHandState.roundState.table)
	g.currentHandState.currentPlayer = p
	g.currentHandState.roundsWon[p]++
}

func (g *Game) isHandOver() bool {
	return g.currentHandState.roundsPlayed == 10
}

func (g *Game) reportSuccess() {
	g.reportSuccessToOwner()

	owner := g.currentHandState.bidWinner
	for p := range g.players {
		if p == owner {
			continue
		}

		if g.currentHandState.roundsWon[p] >= 2 || 10-g.currentHandState.roundsWon[owner] <= 6 {
			g.addResponse(&PlayerResultResponse{Success, p})
		} else {
			p.score.main += int(g.currentHandState.gameType) * 2
			g.addResponse(&PlayerResultResponse{Failiure, p})
		}

		p.score.soups[owner] += g.currentHandState.roundsWon[p] * int(g.currentHandState.gameType) * 2
	}
}

func (g *Game) reportSuccessToOwner() {
	owner := g.currentHandState.bidWinner
	if g.currentHandState.roundsWon[owner] >= 6 {
		owner.score.main -= int(g.currentHandState.gameType) * 2
		g.addResponse(&PlayerResultResponse{Success, owner})
	} else {
		owner.score.main += int(g.currentHandState.gameType) * 2
		g.addResponse(&PlayerResultResponse{Failiure, owner})
	}
}

func (g *Game) addResponse(r Response) {
	g.responses = append(g.responses, r)
}
