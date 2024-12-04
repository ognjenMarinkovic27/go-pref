package game

type PlayerScore struct {
	main  int
	soups map[*Player]int
}

type Player struct {
	pid string

	hand   [10]Card
	played map[Card]bool
	score  PlayerScore

	next *Player
}

func newPlayer(pid string) *Player {
	p := &Player{
		pid: pid,
		score: PlayerScore{
			soups: make(map[*Player]int),
		},
		played: make(map[Card]bool),
	}
	p.next = p
	return p
}
