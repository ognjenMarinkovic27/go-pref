package game

type GameTypeChosenResponse struct {
	GameType GameType
}

func (r *GameTypeChosenResponse) Type() string {
	return "game-type"
}

func (r *GameTypeChosenResponse) RecepientPid() string {
	return ""
}
