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

type HigestAtx struct {
	AtxHex    string `json:"atxHex"`
	AtxBase64 string `json:"atxBase64"`
}

type NetworkInfo struct {
	Epoch             uint32 `json:"epoch"`
	TotalCommited     int64  `json:"totalCommited"`
	CirculatingSupply int64  `json:"circulatingSupply"`
	TotalAccounts     int64  `json:"totalAccounts"`
}
