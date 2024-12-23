package game

type NewBidResponse struct {
	BidderPid string `json:"bidderPid"`
	Bid       Bid    `json:"bid"`
}

func (r *NewBidResponse) Type() string {
	return "new-bid"
}

func (r *NewBidResponse) RecepientPid() string {
	return ""
}
