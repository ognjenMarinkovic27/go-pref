package game

type Response interface {
	Type() string
	RecepientPid() string
}
