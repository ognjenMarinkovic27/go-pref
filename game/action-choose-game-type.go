package game

type ChooseGameTypeAction struct {
	GameType   GameType `json:"game-type"`
	ActionBase `json:"-"`
}

func (action ChooseGameTypeAction) validate(g *Game) bool {
	player := g.players[action.ppid]
	if !g.isCurrentPlayer(player) ||
		g.gameState != ChoosingGameTypeGameState ||
		action.GameType < GameType(g.currentHandState.bid) {
		return false
	}

	return true
}

func (action ChooseGameTypeAction) apply(g *Game) {
	g.chooseGameType(action.GameType)
	g.gameState = RespondingToGameTypeGameState
	g.resetPassed()
	g.moveToNextActivePlayer()
}
