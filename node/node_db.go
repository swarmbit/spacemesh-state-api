package node

import (
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type NodeDB struct {
	DB *sql.Database
}

func NewNodeDB(filePath string, connections int) (*NodeDB, error) {
	db, err := sql.Open(filePath, sql.WithConnections(connections))
	return &NodeDB{
		DB: db,
	}, err
}

func (n *NodeDB) ListRewards(offset int64, limit int64) (rst []*types.NodeSmesherReward, err error) {
	_, err = n.DB.Exec("select r.layer, r.total_reward, r.layer_reward, a.id, a.pubkey, r.coinbase  from rewards_atxs r left join atxs a ON r.atx_id = a.id order by layer limit ?1 offset ?2;",
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, limit)
			stmt.BindInt64(2, offset)

		}, func(stmt *sql.Statement) bool {

			var atxID sTypes.ATXID
			stmt.ColumnBytes(3, atxID[:])
			var nodeID sTypes.NodeID
			stmt.ColumnBytes(4, nodeID[:])
			var address sTypes.Address
			stmt.ColumnBytes(5, address[:])
			reward := &types.NodeSmesherReward{
				Address:     address,
				Layer:       sTypes.LayerID(uint32(stmt.ColumnInt64(0))),
				TotalReward: uint64(stmt.ColumnInt64(1)),
				LayerReward: uint64(stmt.ColumnInt64(2)),
				AtxID:       atxID,
				NodeID:      nodeID,
			}
			rst = append(rst, reward)
			return true
		})
	return
}

func (n *NodeDB) ListLayers(offset int64, limit int64) (layers []*types.NodeLayer, err error) {
	_, err = n.DB.Exec("select id from layers where processed = 1 order by id limit ?1 offset ?2;",
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, limit)
			stmt.BindInt64(2, offset)

		}, func(stmt *sql.Statement) bool {
			layer := &types.NodeLayer{
				Layer: sTypes.LayerID(uint32(stmt.ColumnInt64(0))),
			}
			layers = append(layers, layer)
			return true
		})
	return
}

func (n *NodeDB) ListAccounts(offset int64, limit int64) (accounts []*types.NodeAccount, err error) {
	_, err = n.DB.Exec("select address, layer_updated, next_nonce, balance, template, state from accounts order by layer_updated asc limit ?1 offset ?2;",
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, limit)
			stmt.BindInt64(2, offset)

		}, func(stmt *sql.Statement) bool {
			var atxID sTypes.ATXID
			stmt.ColumnBytes(3, atxID[:])

			var address sTypes.Address
			stmt.ColumnBytes(0, address[:])

			layerUpdated := sTypes.LayerID(uint32(stmt.ColumnInt64(1)))
			nextNonce := stmt.ColumnInt(2)
			balance := stmt.ColumnInt64(3)

			template := stmt.ColumnText(4)

			var state []byte
			stmt.ColumnBytes(5, state[:])

			account := &types.NodeAccount{
				Address:      address,
				LayerUpdated: layerUpdated,
				NextNonce:    nextNonce,
				Balance:      balance,
				Template:     template,
				State:        state,
			}
			accounts = append(accounts, account)
			return true
		})
	return
}

func (n *NodeDB) Close() {
	n.DB.Close()
}
