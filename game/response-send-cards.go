package game

type SendCardsResponse struct {
	PlayerPid string `json:"-"`
}

func (r *SendCardsResponse) Type() string {
	return "send-cards"
}

func (r *SendCardsResponse) RecepientPid() string {
	return r.PlayerPid
}
