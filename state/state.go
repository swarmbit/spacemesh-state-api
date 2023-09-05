package state

import (
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql"
)

type SmesherReward struct {
	Layer       types.LayerID
	TotalReward uint64
	LayerReward uint64
	Coinbase    types.Address
	AtxID       types.ATXID
	NodeID      types.NodeID
}

func CountTotalRewards(db sql.Executor, coinbase types.Address) (count int64, err error) {
	_, err = db.Exec("select count(coinbase) from rewards_atxs where coinbase = ?1;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])

		}, func(stmt *sql.Statement) bool {
			count = stmt.ColumnInt64(0)
			return true
		})
	return
}

// List rewards from all layers for the coinbase address.
func ListRewardsPaginated(db sql.Executor, coinbase types.Address, offset int64, limit int64) (rst []*SmesherReward, err error) {
	_, err = db.Exec("select r.layer, r.total_reward, r.layer_reward, a.id, a.pubkey  from rewards_atxs r left join atxs a ON r.atx_id = a.id where r.coinbase = ?1 order by layer limit ?2 offset ?3;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])
			stmt.BindInt64(2, limit)
			stmt.BindInt64(3, offset)

		}, func(stmt *sql.Statement) bool {

			var atxID types.ATXID
			stmt.ColumnBytes(3, atxID[:])
			var nodeID types.NodeID
			stmt.ColumnBytes(4, nodeID[:])
			reward := &SmesherReward{
				Coinbase:    coinbase,
				Layer:       types.LayerID(uint32(stmt.ColumnInt64(0))),
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

func GetTotalCommittedForEpoch(db sql.Executor, epoch types.EpochID) (sum int64, err error) {
	_, err = db.Exec(`select sum(effective_num_units) from atxs left join identities using(pubkey) 
	where identities.pubkey is null and epoch == ?1;`,
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, int64(epoch))
		}, func(stmt *sql.Statement) bool {
			sum = stmt.ColumnInt64(0)
			return true
		})
	return
}

/*
	db, err := sql.Open("<sql lite file>",
		sql.WithConnections(5),
	)

	if err != nil {
		fmt.Print("Failed to open db")
	}

	address, err := types.StringToAddress("<account>")
	if err != nil {
		fmt.Print("Failed to convert string to address")
	}

	account, err := accounts.Latest(db, address)
	if err != nil {
		fmt.Print("Failed to get account")
	}
	fmt.Printf("Account balance: %d\n", account.Balance)

	rewards, err := ListPaginated(db, address, 0, 1000)
	if err != nil {
		fmt.Printf("Failed to get rewards, error: %s ", err.Error())
	}

	for i, r := range rewards {
		fmt.Printf("reward at index %d: %d, atx id: %s, node id: %s\n", i, r.Layer, r.AtxID.String(), r.NodeID.String())
	}

	totalRewards, err := CountTotalRewards(db, address)
	if err != nil {
		fmt.Printf("Failed to get total rewards, error: %s ", err.Error())
	}
	fmt.Printf("Total rewards: %d\n", totalRewards)

	transactions, err := transactions.GetByAddress(db, types.LayerID(0), types.LayerID(10000), address)
	if err != nil {
		fmt.Printf("Failed to get transactions, error: %s ", err.Error())
	}

	for i, t := range transactions {
		fmt.Printf("transactions at index %d: %d\n", i, t.LayerID)
	}

	highestAtx, _ := atxs.GetIDWithMaxHeight(db, types.EmptyNodeID)
	fmt.Printf("highest: %s\n", highestAtx.String())

	totalForEpoch, _ := GetTotalCommittedForEpoch(db, types.EpochID(3))
	fmt.Printf("totalForEpoch: %f\n", float64(totalForEpoch*64)/float64(1024)/float64(1024))
*/
