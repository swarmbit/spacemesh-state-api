package types

type ActiveNodesEpoch struct {
	Nodes []string `json:"nodes"`
}

type Epoch struct {
	EffectiveUnitsCommited uint64 `json:"effectiveUnitsCommited"`
	EpochSubsidy           uint64 `json:"epochSubsidy"`
	TotalWeight            uint64 `json:"totalWeight"`
	TotalRewards           int64  `json:"totalRewards"`
	TotalActiveSmeshers    uint64 `json:"totalActiveSmeshers"`
}

type Atx struct {
	NodeId            string `json:"nodeId"`
	AtxId             string `json:"atxId"`
	EffectiveNumUnits uint32 `json:"effectiveNumUnits"`
	Received          int64  `json:"received"`
}

type ShortAccount struct {
	Balance  uint64 `json:"balance"`
	USDValue int64  `json:"usdValue"`
	Address  string `json:"address"`
}

type AccountPostResponse struct {
	Account                string `json:"account"`
	TotalEffectiveNumUnits int64  `json:"totalEffectiveNumUnits"`
}

type Account struct {
	Balance              uint64 `json:"balance"`
	USDValue             int64  `json:"usdValue"`
	BalanceDisplay       string `json:"balanceDisplay"`
	NumberOfTransactions int64  `json:"numberOfTransactions"`
	Counter              int64  `json:"counter"`
	NumberOfRewards      int64  `json:"numberOfRewards"`
	TotalRewards         uint64 `json:"totalRewards"`
	Address              string `json:"address"`
}

type Reward struct {
	Account        string `json:"account"`
	Rewards        int64  `json:"rewards"`
	RewardsDisplay string `json:"rewardsDisplay"`
	Layer          int64  `json:"layer"`
	SmesherId      string `json:"smesherId"`
	Time           string `json:"time"`
	Timestamp      int64  `json:"timestamp"`
}

type Transaction struct {
	ID               string `json:"id"`
	Status           uint8  `json:"status"`
	PrincipalAccount string `json:"principalAccount"`
	ReceiverAccount  string `json:"receiverAccount"`
	Fee              uint64 `json:"fee"`
	Amount           uint64 `json:"amount"`
	Layer            uint32 `json:"layer"`
	Counter          uint64 `json:"counter"`
	Method           string `json:"method"`
	Timestamp        int64  `json:"timestamp"`
}

type RewardDetails struct {
	TotalSum                 int64        `json:"totalSum"`
	CurrentEpoch             int64        `json:"currentEpoch"`
	CurrentEpochRewardsSum   int64        `json:"currentEpochRewardsSum"`
	CurrentEpochRewardsCount int64        `json:"currentEpochRewardsCount"`
	Eligibility              *Eligibility `json:"eligibility"`
}

type RewardDetailsEpoch struct {
	Epoch        int64        `json:"epoch"`
	RewardsSum   int64        `json:"rewardsSum"`
	RewardsCount int64        `json:"rewardsCount"`
	Eligibility  *Eligibility `json:"eligibility"`
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
