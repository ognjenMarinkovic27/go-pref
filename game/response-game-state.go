package game

type GameStateResponse struct {
	GameState GameState
}

func (r *GameStateResponse) Type() string {
	return "game-state"
}

func (r *GameStateResponse) RecepientPid() string {
	return ""
}
