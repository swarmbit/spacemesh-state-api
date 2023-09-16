package types

type TrackingDoc struct {
	Id     string `bson:"_id"`
	Offset int64  `bson:"offset"`
}

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

type AccountDoc struct {
	Address      string `bson:"_id"`
	LayerUpdated int64  `bson:"layer_updated"`
	Balance      int64  `bson:"balance"`
	NextNonce    int    `bson:"next_nonce"`
	Template     string `bson:"template"`
	State        []byte `bson:"state"`
}
