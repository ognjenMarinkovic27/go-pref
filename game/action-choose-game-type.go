package game

type ChooseGameTypeAction struct {
	gameType GameType
	player   *Player
}

func NewChooseGameTypeAction(gameType GameType, player *Player) ChooseGameTypeAction {
	return ChooseGameTypeAction{gameType, player}
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
