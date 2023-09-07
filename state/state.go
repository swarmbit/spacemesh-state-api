package state

import (
	"fmt"

	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/proposals/util"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/spacemeshos/go-spacemesh/sql/atxs"
	"github.com/spacemeshos/go-spacemesh/tortoise"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type State struct {
	DB sql.Executor
}

func NewState() *State {
	db, err := StartDB("/Users/brunovale/dev/git/spacemesh/spacemesh-configs/custom-node/node-data/state.sql", 10)
	if err != nil {
		fmt.Print("Failed to open db")
	}
	return &State{
		DB: db,
	}
}

func (s *State) CountTotalRewards(coinbase sTypes.Address) (count int64, err error) {
	_, err = s.DB.Exec("select count(coinbase) from rewards_atxs where coinbase = ?1;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])

		}, func(stmt *sql.Statement) bool {
			count = stmt.ColumnInt64(0)
			return true
		})
	return
}

func getTotalWeight(db sql.Executor, epoch sTypes.EpochID) (total uint64, err error) {
	_, err = db.Exec("select tick_count, effective_num_units from atxs where epoch = ?1;",
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, int64(epoch))
		}, func(stmt *sql.Statement) bool {
			tickCount := uint64(stmt.ColumnInt64(0))
			effectiveNumUnits := uint64(stmt.ColumnInt64(1))
			total += GetATXWeight(effectiveNumUnits, tickCount)
			return true
		})
	return
}

func (s *State) GetLatestAVGLayerReward() (uint64, error) {
	var sumRewards uint64
	_, err := s.DB.Exec("select DISTINCT(layer), layer_reward from rewards_atxs ORDER By layer DESC LIMIT 2000;",
		func(stmt *sql.Statement) {
		}, func(stmt *sql.Statement) bool {
			sumRewards += uint64(stmt.ColumnInt64(1))
			return true
		})
	return sumRewards / 2000, err
}

// List rewards from all layers for the coinbase address.
func (s *State) ListRewardsPaginated(coinbase sTypes.Address, offset int64, limit int64) (rst []*types.SmesherReward, err error) {
	_, err = s.DB.Exec("select r.layer, r.total_reward, r.layer_reward, a.id, a.pubkey  from rewards_atxs r left join atxs a ON r.atx_id = a.id where r.coinbase = ?1 order by layer limit ?2 offset ?3;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])
			stmt.BindInt64(2, limit)
			stmt.BindInt64(3, offset)

		}, func(stmt *sql.Statement) bool {

			var atxID sTypes.ATXID
			stmt.ColumnBytes(3, atxID[:])
			var nodeID sTypes.NodeID
			stmt.ColumnBytes(4, nodeID[:])
			reward := &types.SmesherReward{
				Coinbase:    coinbase,
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

func GetATXWeight(numUnits, tickCount uint64) uint64 {
	return safeMul(numUnits, tickCount)
}

func safeMul(a, b uint64) uint64 {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		panic("uint64 overflow")
	}
	return c
}

func (s *State) GetEligibilityCount(nodeID sTypes.NodeID, epoch sTypes.EpochID) (uint32, error) {

	defaultConfig := tortoise.DefaultConfig()
	layerSize := defaultConfig.LayerSize
	minimalWeight := defaultConfig.MinimalActiveSetWeight

	atx, err := atxs.GetByEpochAndNodeID(s.DB, epoch, nodeID)
	if err != nil {
		fmt.Printf("Failed to get atx for node, error: %s ", err.Error())
		return 0, err
	}

	atxWeight := atx.GetWeight()

	total, err := getTotalWeight(s.DB, epoch)
	if err != nil {
		fmt.Printf("Failed to get total weight, error: %s ", err.Error())
		return 0, err
	}

	slots, err := util.GetNumEligibleSlots(atxWeight, minimalWeight, total, layerSize, 4032)
	return uint32(slots), err
}

func (s *State) GetTotalCommittedForEpoch(epoch sTypes.EpochID) (sum int64, err error) {
	_, err = s.DB.Exec(`select sum(effective_num_units) from atxs left join identities using(pubkey) 
	where identities.pubkey is null and epoch == ?1;`,
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, int64(epoch))
		}, func(stmt *sql.Statement) bool {
			sum = stmt.ColumnInt64(0)
			return true
		})
	return
}
