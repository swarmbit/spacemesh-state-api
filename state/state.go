package state

import (
	"fmt"
	"sync"
	"time"

	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/proposals/util"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/spacemeshos/go-spacemesh/sql/atxs"
	"github.com/spacemeshos/go-spacemesh/sql/layers"
	"github.com/spacemeshos/go-spacemesh/tortoise"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type State struct {
	DocDB             *DocDB
	DB                sql.Executor
	DBInstance        *sql.Database
	epochsTotalWeight *sync.Map
	avgRewards        *sync.Map
}

func NewState() *State {
	db, err := StartDB("/Users/brunovale/spacemesh-db/node-data/state.sql", 1)
	if err != nil {
		fmt.Print("Failed to open db")
	}

	docDB, _ := NewDocDB()
	stateObj := &State{
		DocDB:             docDB,
		DB:                db,
		DBInstance:        db,
		epochsTotalWeight: &sync.Map{},
		avgRewards:        &sync.Map{},
	}
	stateObj.periodicDocDB()
	return stateObj
}

func (s *State) periodicDocDB() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			offset, _ := s.DocDB.GetOffset("rewards")
			fmt.Println("Next rewards offset", offset)
			rewards, _ := s.ListRewardsNoAddress(offset, 100)
			rewardsDoc := make([]interface{}, len(rewards))
			for i, r := range rewards {
				rewardsDoc[i] = RewardsDoc{
					Ammount: int64(r.TotalReward),
					AtxID:   r.AtxID.String(),
					Layer:   int64(r.Layer),
				}
			}
			s.DocDB.SaveRewards(offset, rewardsDoc)
		}
	}()
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

func (s *State) CountTotalRewardsEpoch(coinbase sTypes.Address, epoch sTypes.EpochID) (count int64, err error) {
	firstLayer := uint64(epoch) * config.LayersPerEpoch
	lastLayer := firstLayer + config.LayersPerEpoch - 1
	_, err = s.DB.Exec("select count(coinbase) from rewards_atxs where coinbase = ?1 AND layer >= ?2 AND layer <= ?3;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])
			stmt.BindInt64(2, int64(firstLayer))
			stmt.BindInt64(3, int64(lastLayer))
		}, func(stmt *sql.Statement) bool {
			count = stmt.ColumnInt64(0)
			return true
		})
	return
}

func (s *State) SumRewards(coinbase sTypes.Address) (count int64, err error) {
	_, err = s.DB.Exec("select sum(total_reward) from rewards_atxs where coinbase = ?1;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])

		}, func(stmt *sql.Statement) bool {
			count = stmt.ColumnInt64(0)
			return true
		})
	return
}

func (s *State) SumRewardsEpoch(coinbase sTypes.Address, epoch sTypes.EpochID) (count int64, err error) {
	firstLayer := uint64(epoch) * config.LayersPerEpoch
	lastLayer := firstLayer + config.LayersPerEpoch - 1
	_, err = s.DB.Exec("select sum(total_reward) from rewards_atxs where coinbase = ?1 AND layer >= ?2 AND layer <= ?3;",
		func(stmt *sql.Statement) {
			stmt.BindBytes(1, coinbase[:])
			stmt.BindInt64(2, int64(firstLayer))
			stmt.BindInt64(3, int64(lastLayer))
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

func (s *State) GetCurrentEpoch() (sTypes.EpochID, sTypes.LayerID, error) {
	currentLayer, err := layers.GetProcessed(s.DB)
	if err != nil {
		fmt.Printf("Failed to get current layer, error: %s ", err.Error())
		return sTypes.EpochID(0), currentLayer, err
	}
	return GetEpoch(currentLayer), currentLayer, err
}

func (s *State) GetLatestAVGLayerReward() (uint64, error) {

	var rewardsAvg uint64
	currentLayer, err := layers.GetProcessed(s.DB)
	if err != nil {
		fmt.Printf("Failed to get current layer, error: %s ", err.Error())
		return 0, err
	}

	avgCached, exist := s.avgRewards.Load(currentLayer.Uint32())
	if !exist {
		var sumRewards uint64
		_, err := s.DB.Exec("select DISTINCT(layer), layer_reward from rewards_atxs ORDER By layer DESC LIMIT 2000;",
			func(stmt *sql.Statement) {
			}, func(stmt *sql.Statement) bool {
				sumRewards += uint64(stmt.ColumnInt64(1))
				return true
			})
		if err != nil {
			fmt.Printf("Failed to calculate AVG, error: %s ", err.Error())
			return 0, err
		}
		rewardsAvg = sumRewards / 2000
		result, _ := s.avgRewards.LoadOrStore(currentLayer.Uint32(), rewardsAvg)
		avgCached = result.(uint64)
	}

	return avgCached.(uint64), nil
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

func (s *State) ListRewardsNoAddress(offset int64, limit int64) (rst []*types.SmesherReward, err error) {
	_, err = s.DB.Exec("select r.layer, r.total_reward, r.layer_reward, a.id, a.pubkey  from rewards_atxs r left join atxs a ON r.atx_id = a.id order by layer limit ?1 offset ?2;",
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, limit)
			stmt.BindInt64(2, offset)

		}, func(stmt *sql.Statement) bool {

			var atxID sTypes.ATXID
			stmt.ColumnBytes(3, atxID[:])
			var nodeID sTypes.NodeID
			stmt.ColumnBytes(4, nodeID[:])
			reward := &types.SmesherReward{
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

func (s *State) GetEligibilityCount(nodeID sTypes.NodeID) (int32, error) {

	defaultConfig := tortoise.DefaultConfig()
	layerSize := defaultConfig.LayerSize
	minimalWeight := defaultConfig.MinimalActiveSetWeight

	currentLayer, err := layers.GetProcessed(s.DB)
	if err != nil {
		fmt.Printf("Failed to get current layer, error: %s ", err.Error())
		return 0, err
	}

	epoch := GetEpoch(currentLayer) - 1

	atx, err := atxs.GetByEpochAndNodeID(s.DB, epoch, nodeID)
	if err != nil {
		fmt.Printf("Failed to get atx for node, error: %s ", err.Error())
		if atx == nil {
			return -1, err
		}
		return 0, err
	}

	atxWeight := atx.GetWeight()

	totalCached, exist := s.epochsTotalWeight.Load(epoch.Uint32())
	if !exist {
		println("Load total cached")
		total, err := getTotalWeight(s.DB, epoch)
		if err != nil {
			fmt.Printf("Failed to get total weight, error: %s ", err.Error())
			return 0, err
		}
		totalCached, _ = s.epochsTotalWeight.LoadOrStore(epoch.Uint32(), total)
	}

	slots, err := util.GetNumEligibleSlots(atxWeight, minimalWeight, totalCached.(uint64), layerSize, 4032)
	return int32(slots), err
}

func GetEpoch(l sTypes.LayerID) sTypes.EpochID {
	return sTypes.EpochID(l.Uint32() / config.LayersPerEpoch)
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

func (s *State) GetCirculatingSupply() (int64, error) {
	var circulatingSupply int64
	_, err := s.DB.Exec(`WITH LatestBalanceChanges AS (
				SELECT address, MAX(layer_updated) as MaxLayerUpdated
				FROM accounts
				GROUP BY address
			),
			LatestBalances AS (
				SELECT a.address, a.balance
				FROM accounts a
				JOIN LatestBalanceChanges l ON a.address = l.address AND a.layer_updated = l.MaxLayerUpdated
			)
			SELECT SUM(balance) as TotalLatestBalance
			FROM LatestBalances;`,
		func(stmt *sql.Statement) {
		}, func(stmt *sql.Statement) bool {
			circulatingSupply = stmt.ColumnInt64(0)
			return true
		})
	if err != nil {
		fmt.Printf("Failed to get total balances, error: %s ", err.Error())
		return 0, err
	}
	return circulatingSupply - 150000000000000000, nil
}

func (s *State) GetNumberOfAccounts() (sum int64, err error) {
	_, err = s.DB.Exec(`select count(DISTINCT(address)) from accounts;`,
		func(stmt *sql.Statement) {
		}, func(stmt *sql.Statement) bool {
			sum = stmt.ColumnInt64(0)
			return true
		})
	return
}

func (s *State) GetActiveAtxCount(epoch sTypes.EpochID) (count int64, err error) {
	_, err = s.DB.Exec(`select count(id) from atxs where epoch = ?1;`,
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, int64(epoch))
		}, func(stmt *sql.Statement) bool {
			count = stmt.ColumnInt64(0)
			return true
		})
	return
}

func (s *State) GetActiveAtxPerAddress(epoch sTypes.EpochID, coinbase sTypes.Address) (smeshers []*types.Smesher, err error) {
	_, err = s.DB.Exec(`select pubkey, effective_num_units from atxs where epoch = ?1 and coinbase = ?2;`,
		func(stmt *sql.Statement) {
			stmt.BindInt64(1, int64(epoch))
			stmt.BindBytes(2, coinbase[:])
		}, func(stmt *sql.Statement) bool {
			var nodeID sTypes.NodeID
			stmt.ColumnBytes(0, nodeID[:])
			numUnits := stmt.ColumnInt64(1)
			smeshers = append(smeshers, &types.Smesher{
				NodeID:            nodeID,
				Coinbase:          coinbase,
				EffectiveNumUnits: numUnits,
			})
			return true
		})
	return
}

func (s *State) GetHighestAtx() (sTypes.ATXID, error) {
	return atxs.GetIDWithMaxHeight(s.DB, sTypes.EmptyNodeID)
}
