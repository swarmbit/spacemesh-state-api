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

type RewardDetails struct {
	TotalSum                 int64 `json:"totalSum"`
	CurrentEpoch             int64 `json:"currentEpoch"`
	CurrentEpochRewardsSum   int64 `json:"currentEpochRewardsSum"`
	CurrentEpochRewardsCount int64 `json:"currentEpochRewardsCount"`
}

type Eligibility struct {
	Count             int32  `json:"count"`
	SmesherId         string `json:"smesherId"`
	EffectiveNumUnits int64  `json:"effectiveNumUnits"`
	PredictedRewards  uint64 `json:"predictedRewards"`
}

type NetworkInfo struct {
	Epoch                  uint32 `json:"epoch"`
	Layer                  int64  `json:"layer"`
	EffectiveUnitsCommited int64  `json:"effectiveUnitsCommited"`
	CirculatingSupply      int64  `json:"circulatingSupply"`
	TotalAccounts          int64  `json:"totalAccounts"`
	TotalActiveSmeshers    int64  `json:"totalActiveSmeshers"`
	AtxHex                 string `json:"atxHex"`
	AtxBase64              string `json:"atxBase64"`
}
