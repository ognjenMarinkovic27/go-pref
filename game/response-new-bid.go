package game

type NewBidResponse struct {
	BidderPid string
	Bid       Bid
}

func (r *NewBidResponse) Type() string {
	return "new-bid"
}

func (r *NewBidResponse) RecepientPid() string {
	return ""
}
