package game

type Action interface {
	validate(g *Game) bool
	apply(g *Game)

	/* TODO: stinky? */
	playerPid() string
	/* TODO: stinky? */
	SetPlayerPid(ppid string)
}

type ActionBase struct {
	ppid string
}

func (a *ActionBase) playerPid() string {
	return a.ppid
}

func (a *ActionBase) SetPlayerPid(ppid string) {
	a.ppid = ppid
}
