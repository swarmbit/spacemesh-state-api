package processor

import (
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/config"
)

type SyncProcessor struct {
	WriteDB                *database.WriteDB
}

func NewSyncProcessor(configValues *config.Config, writeDB *database.WriteDB) *SyncProcessor {
	return &SyncProcessor{
		WriteDB:                writeDB,
	}
}

func (*SyncProcessor) StartProcessingSync()  {

}
