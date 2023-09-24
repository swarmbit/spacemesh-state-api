package types

type RewardsDoc struct {
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

type AtxDoc struct {
	AtxID             string `bson:"_id"`
	NodeID            string `bson:"node_id"`
	Coinbase          string `bson:"coinbase"`
	PublishEpoch      uint32 `json:"publish_epoch"`
	EffectiveNumUnits uint32 `bson:"effective_num_units"`
	BaseTick          uint64 `bson:"base_tick"`
	TickCount         uint64 `bson:"tick_count"`
	Sequence          uint64 `json:"sequence"`
	Received          int64  `json:"received"`
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
}

type AccountDoc struct {
	Address      string `bson:"_id"`
	LayerUpdated int64  `bson:"layer_updated"`
	Balance      int64  `bson:"balance"`
	NextNonce    int    `bson:"next_nonce"`
	Template     string `bson:"template"`
	State        []byte `bson:"state"`
}
