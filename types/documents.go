package types

type RewardsDoc struct {
	Id          string `bson:"_id"`
	NodeId      string `bson:"node_id"`
	Coinbase    string `bson:"coinbase"`
	AtxID       string `bson:"atx_id"`
	LayerReward int64  `bson:"layerReward"`
	TotalReward int64  `bson:"totalReward"`
	Layer       int64  `bson:"layer"`
}

type LayerDoc struct {
	Layer  int64 `bson:"_id"`
	Status int   `bson:"status"`
}

type NodeDoc struct {
	ID          string             `bson:"_id"`
	Atxs        []NodeAtxDoc       `bson:"atxs"`
	Malfeasance MalfeasanceNodeDoc `bson:"malfeasance"`
}

type MalfeasanceNodeDoc struct {
	Received int64 `json:"received"`
}

type NodeAtxDoc struct {
	Coinbase          string `bson:"coinbase"`
	PublishEpoch      uint32 `json:"publish_epoch"`
	EffectiveNumUnits uint32 `bson:"effectiveNumUnits"`
	Weight            uint64 `bson:"weight"`
	Sequence          uint64 `json:"sequence"`
	Received          int64  `json:"received"`
}

type AtxDoc struct {
	AtxID             string `bson:"_id"`
	NodeID            string `bson:"node_id"`
	Coinbase          string `bson:"coinbase"`
	PublishEpoch      uint32 `json:"publish_epoch"`
	EffectiveNumUnits uint32 `bson:"effective_num_units"`
	BaseTick          uint64 `bson:"base_tick"`
	Weight            uint64 `bson:"weight"`
	TickCount         uint64 `bson:"tick_count"`
	Sequence          uint64 `json:"sequence"`
	Received          int64  `json:"received"`
}

type AtxEpochDoc struct {
	ID                     int64  `bson:"_id"`
	TotalEffectiveNumUnits uint64 `bson:"totalEffectiveNumUnits"`
	TotalWeight            uint64 `bson:"totalWeight"`
}

type TransactionDoc struct {
	ID              string `bson:"_id"`
	Status          uint8  `json:"status"`
	PrincipaAccount string `bson:"principal_account"`
	ReceiverAccount string `bson:"receiver_account"`
	Fee             uint64 `bson:"fee"`
	Gas             uint64 `bson:"gas"`
	GasPrice        uint64 `bson:"gas_price"`
	Amount          uint64 `bson:"amount"`
	Layer           uint32 `bson:"layer"`
	Counter         uint64 `bson:"counter"`
	Method          uint8  `json:"method"`
	Complete        bool   `json:"complete"`
}

type AccountDoc struct {
	Address      string `bson:"_id"`
	Balance      uint64 `bson:"balance"`
	TotalRewards uint64 `bson:"totalRewards"`
	Fees         uint64 `bson:"fees"`
	Sent         uint64 `bson:"sent"`
}

type NetworkInfoDoc struct {
	Id                string `bson:"_id"`
	CirculatingSupply uint64 `bson:"circulatingSupply"`
}

type AccountPost struct {
	Id                     *AccountPostId `bson:"_id"`
	TotalEffectiveNumUnits int64          `bson:"totalEffectiveNumUnits"`
}

type AccountPostId struct {
	Coinbase string `bson:"coinbase"`
}
type AggregationTotal struct {
	TotalSum int64 `bson:"totalSum"`
}

type AggregationAtxTotals struct {
	TotalWeight            int64 `bson:"totalWeight"`
	TotalEffectiveNumUnits int64 `bson:"totalEffectiveNumUnits"`
}
