package types

type Malfeasance struct {
	LayerID  uint32 `json:"layer"`
	NodeID   string `json:"node_id"`
	Received int64  `json:"received"`
}

type LayerUpdate struct {
	LayerID uint32 `json:"layer"`
	Status  int    `json:"status"`
}

type Reward struct {
	ID          string `json:"id"`
	Layer       uint32 `json:"layer"`
	Total       uint64 `json:"totalReward"`
	LayerReward uint64 `json:"layerReward"`
	Coinbase    string `json:"coinbase"`
	AtxID       string `json:"atxID"`
	NodeID      string `json:"nodeID"`
}

type Atx struct {
	Received          int64  `json:"received"`
	BaseTick          uint64 `json:"baseTick"`
	TickCount         uint64 `json:"tickCount"`
	EffectiveNumUnits uint32 `json:"EffectiveNumUnits"`
	AtxID             string `json:"atxID"`
	NodeID            string `json:"nodeID"`
	Sequence          uint64 `json:"sequence"`
	PublishEpoch      uint32 `json:"publishEpoch"`
	Coinbase          string `json:"coinbase"`
}

type Transaction struct {
	ID     string             `json:"id"`
	Header *TransactionHeader `json:"header"`
	Raw    []byte             `json:"raw"`
}

type TransactionHeader struct {
	Message         string   `json:"message"`
	Status          uint8    `json:"status"`
	BlockID         string   `json:"block_id"`
	LayerID         uint32   `json:"layer_id"`
	Principal       string   `json:"principal"`
	TemplateAddress string   `json:"template_address"`
	Method          uint8    `json:"method"`
	Nonce           uint64   `json:"nonce"`
	Gas             uint64   `json:"gas"`
	Fee             uint64   `json:"fee"`
	Addresses       []string `json:"addresses"`
}