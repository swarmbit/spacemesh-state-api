package types

type Account struct {
	Balance              uint64 `json:"balance"`
	BalanceDisplay       string `json:"balanceDisplay"`
	NumberOfTransactions int64  `json:"numberOfTransactions"`
	NumberOfRewards      int64  `json:"numberOfRewards"`
	TotalRewards         uint64 `json:"totalRewards"`
	Address              string `json:"address"`
}

type Reward struct {
	Rewards        int64  `json:"rewards"`
	RewardsDisplay string `json:"rewardsDisplay"`
	Layer          int64  `json:"layer"`
	SmesherId      string `json:"smesherId"`
	Time           string `json:"time"`
	Timestamp      int64  `json:"timestamp"`
}

type Transaction struct {
	ID              string `bson:"_id"`
	Status          uint8  `json:"status"`
	PrincipaAccount string `bson:"principal_account"`
	ReceiverAccount string `bson:"receiver_account"`
	Fee             uint64 `bson:"fee"`
	Amount          uint64 `bson:"amount"`
	Layer           uint32 `bson:"layer"`
	Counter         uint64 `bson:"counter"`
	Method          string `json:"method"`
}

type RewardDetails struct {
	TotalSum                 int64 `json:"totalSum"`
	CurrentEpoch             int64 `json:"currentEpoch"`
	CurrentEpochRewardsSum   int64 `json:"currentEpochRewardsSum"`
	CurrentEpochRewardsCount int64 `json:"currentEpochRewardsCount"`
}

type Eligibility struct {
	Address           string `json:"address"`
	Count             int32  `json:"count"`
	EffectiveNumUnits int64  `json:"effectiveNumUnits"`
	PredictedRewards  uint64 `json:"predictedRewards"`
}

type NetworkInfo struct {
	Epoch                  uint32                `json:"epoch"`
	Layer                  uint64                `json:"layer"`
	EffectiveUnitsCommited uint64                `json:"effectiveUnitsCommited"`
	CirculatingSupply      uint64                `json:"circulatingSupply"`
	TotalAccounts          uint64                `json:"totalAccounts"`
	TotalActiveSmeshers    uint64                `json:"totalActiveSmeshers"`
	AtxHex                 string                `json:"atxHex"`
	AtxBase64              string                `json:"atxBase64"`
	NextEpoch              *NetworkInfoNextEpoch `json:"nextEpoch"`
}

type NetworkInfoNextEpoch struct {
	Epoch                  uint32 `json:"epoch"`
	EffectiveUnitsCommited int64  `json:"effectiveUnitsCommited"`
	TotalActiveSmeshers    int64  `json:"totalActiveSmeshers"`
}
