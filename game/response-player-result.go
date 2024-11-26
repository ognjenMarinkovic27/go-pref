package game

const (
	Success  = true
	Failiure = false
)

type PlayerResultResponse struct {
	Success bool
	Player  *Player
}
