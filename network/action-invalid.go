package network

type InvalidAction struct{}

func (p *InvalidAction) Type() string {
	return "invalid-action"
}
