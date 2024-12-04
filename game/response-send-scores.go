package game

type SendScoresResponse struct{}

func (r *SendScoresResponse) Type() string {
	return "send-scores"
}

func (r *SendScoresResponse) RecepientPid() string {
	return ""
}
