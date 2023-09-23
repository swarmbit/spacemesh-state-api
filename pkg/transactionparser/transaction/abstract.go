package transaction

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
)

// DecodedTransactioner is an interface for transaction decoded from raw bytes.
type DecodedTransactioner interface {
	GetType() uint8
	GetAmount() uint64
	GetCounter() uint64
	GetReceiver() types.Address
	GetGasPrice() uint64
	GetPrincipal() types.Address
	GetPublicKeys() [][]byte
	GetSignature() []byte
}
