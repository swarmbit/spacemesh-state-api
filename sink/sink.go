package sink

import (
	"fmt"
	"time"

	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/node"
)

type Sink struct {
	DocDB  *database.DocDB
	NodeDB *node.NodeDB
}

func NewSink(docDB *database.DocDB, nodeDB *node.NodeDB) *Sink {
	return &Sink{
		DocDB:  docDB,
		NodeDB: nodeDB,
	}
}

func (s *Sink) StartRewardsSink() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			offset, _ := s.DocDB.GetOffset("rewards")
			fmt.Println("Next rewards offset", offset)
			rewards, _ := s.NodeDB.ListRewards(offset, 1000)
			s.DocDB.SaveRewards(offset, rewards)
		}
	}()
}

func (s *Sink) StartLayersSink() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			offset, _ := s.DocDB.GetOffset("layers")
			fmt.Println("Next layers offset", offset)
			layers, _ := s.NodeDB.ListLayers(offset, 1000)
			s.DocDB.SaveLayers(offset, layers)
		}
	}()
}

func (s *Sink) StartAccountsSink() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			offset, _ := s.DocDB.GetOffset("accounts")
			fmt.Println("Next accounts offset", offset)
			accounts, _ := s.NodeDB.ListAccounts(offset, 1000)
			s.DocDB.SaveAccounts(offset, accounts)
		}
	}()
}
