package game

const (
	Coming    = true
	NotComing = false
)

type PlayerComingResponse struct {
	Coming bool
	Player *Player
}
