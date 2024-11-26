package game

type PlayerScore struct {
	score int
	soups map[*Player]int
}

type Player struct {
	hand   [10]Card
	played map[Card]bool
	score  PlayerScore

	next *Player
}

func NewPlayer() *Player {
	p := &Player{
		score: PlayerScore{
			soups: make(map[*Player]int),
		},
		played: make(map[Card]bool),
	}
	p.next = p
	return p
}
