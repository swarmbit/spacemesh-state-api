package types

type Account struct {
	Balance              uint64 `json:"balance"`
	USDValue             uint64 `json:"usdValue"`
	BalanceDisplay       string `json:"balanceDisplay"`
	NumberOfTransactions int64  `json:"numberOfTransactions"`
	Counter              int64  `json:"counter"`
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
	TotalSum                 int64        `json:"totalSum"`
	CurrentEpoch             int64        `json:"currentEpoch"`
	CurrentEpochRewardsSum   int64        `json:"currentEpochRewardsSum"`
	CurrentEpochRewardsCount int64        `json:"currentEpochRewardsCount"`
	Eligibility              *Eligibility `json:"eligibility"`
}

type Eligibility struct {
	Count             int32  `json:"count"`
	EffectiveNumUnits int64  `json:"effectiveNumUnits"`
	PredictedRewards  uint64 `json:"predictedRewards"`
}

type NetworkInfo struct {
	Epoch                  uint32                `json:"epoch"`
	Layer                  uint64                `json:"layer"`
	EffectiveUnitsCommited uint64                `json:"effectiveUnitsCommited"`
	EpochSubsidy           uint64                `json:"epochSubsidy"`
	TotalSlots             uint64                `json:"totalSlots"`
	TotalWeight            uint64                `json:"totalWeight"`
	CirculatingSupply      uint64                `json:"circulatingSupply"`
	Price                  float64               `json:"price"`
	MarketCap              uint64                `json:"marketCap"`
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
