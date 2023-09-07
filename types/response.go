package types

type Account struct {
	Balance        int64  `json:"balance"`
	BalanceDisplay string `json:"balanceDisplay"`
	Counter        int64  `json:"counter"`
	Address        string `json:"address"`
}

type Reward struct {
	Rewards        int64  `json:"rewards"`
	RewardsDisplay string `json:"rewardsDisplay"`
	Layer          int64  `json:"layer"`
	SmesherId      string `json:"smesherId"`
	Time           string `json:"time"`
	Timestamp      int64  `json:"timestamp"`
}

type Eligibility struct {
	Count            int32  `json:"count"`
	SmesherId        string `json:"smesherId"`
	PredictedRewards uint64 `json:"predictedRewards"`
}
