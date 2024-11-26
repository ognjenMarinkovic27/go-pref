package game

type Action interface {
	validate(g *Game) bool
	apply(g *Game)
}
