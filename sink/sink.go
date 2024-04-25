package sink

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	natsS "github.com/spacemeshos/go-spacemesh/nats"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/config"
)

type Sink struct {
	WriteDB                *database.WriteDB
	layersSub              *nats.Subscription
	rewardsSub             *nats.Subscription
	atxSub                 *nats.Subscription
	transactionsResultSub  *nats.Subscription
	transactionsCreatedSub *nats.Subscription
	malfeasanceSub         *nats.Subscription
}

func NewSink(configValues *config.Config, writeDB *database.WriteDB) *Sink {
	nc, err := nats.Connect(configValues.Nats.Uri)
	if err != nil {
		panic("Failed to connect to NATS")
	}
	js, _ := nc.JetStream()

	js.AddConsumer("layers", &nats.ConsumerConfig{
		Durable:        "state-api-process-layers",
		DeliverSubject: "layers",
		DeliverGroup:   "state-api-process-layers",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("rewards", &nats.ConsumerConfig{
		Durable:        "state-api-process-rewards",
		DeliverSubject: "rewards",
		DeliverGroup:   "state-api-process-rewards",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("atx", &nats.ConsumerConfig{
		Durable:        "state-api-process-atx",
		DeliverSubject: "atx",
		DeliverGroup:   "state-api-process-atx",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("transactions", &nats.ConsumerConfig{
		Durable:        "state-api-process-transactions-result",
		DeliverSubject: "transactions.result",
		DeliverGroup:   "state-api-process-transactions",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("transactions", &nats.ConsumerConfig{
		Durable:        "state-api-process-transactions-created",
		DeliverSubject: "transactions.created",
		DeliverGroup:   "state-api-process-transactions",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("malfeasance", &nats.ConsumerConfig{
		Durable:        "state-api-process-malfeasance",
		DeliverSubject: "malfeasance",
		DeliverGroup:   "state-api-process-malfeasance",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	fmt.Println("Connect to nats stream")
	layersSub, err := js.PullSubscribe("layers", "state-api-process-layers", nats.BindStream("layers"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	rewardsSub, err := js.PullSubscribe("rewards", "state-api-process-rewards", nats.BindStream("rewards"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	atxSub, err := js.PullSubscribe("atx", "state-api-process-atx", nats.BindStream("atx"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	transactionsResultSub, err := js.PullSubscribe("transactions.result", "state-api-process-transactions-result", nats.BindStream("transactions"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	transactionsCreatedSub, err := js.PullSubscribe("transactions.created", "state-api-process-transactions-created", nats.BindStream("transactions"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	malfeasanceSub, err := js.PullSubscribe("malfeasance", "state-api-process-malfeasance", nats.BindStream("malfeasance"))
	if err != nil {
		fmt.Println("Failed to subscribe: ", err)
	}
	return &Sink{
		layersSub:              layersSub,
		rewardsSub:             rewardsSub,
		atxSub:                 atxSub,
		transactionsResultSub:  transactionsResultSub,
		transactionsCreatedSub: transactionsCreatedSub,
		malfeasanceSub:         malfeasanceSub,
		WriteDB:                writeDB,
	}
}

func (s *Sink) StartRewardsSink() {
	fmt.Println("Start rewards sink")
	go func() {
		for {
			msgs, err := s.rewardsSub.Fetch(100, nats.MaxWait(2*time.Hour))
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}
			var wg sync.WaitGroup
			for _, msg := range msgs {
				wg.Add(1)
				go s.processRewardMessage(msg, &wg)
			}
			wg.Wait()
		}
	}()
}

func (s *Sink) processRewardMessage(msg *nats.Msg, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("New reward")
	var reward *natsS.Reward
	errJson := json.Unmarshal(msg.Data, &reward)
	fmt.Println("Next reward: ", reward.Layer)
	if errJson != nil {
		fmt.Println("Error parsing json reward: ", errJson)
		msg.Nak()
		return
	}
	saveErr := s.WriteDB.SaveReward(reward)

	if saveErr != nil {
		fmt.Println("Failed to save reward")
		msg.Nak()
	} else {
		fmt.Println("Reward saved")
		msg.AckSync()
	}
}

func (s *Sink) StartLayersSink() {
	fmt.Println("Start layers sink")

	go func() {
		for {
			msgs, err := s.layersSub.Fetch(100, nats.MaxWait(2*time.Hour))
			fmt.Println("New layers")
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}
			for _, msg := range msgs {
				fmt.Println("Layer: ", string(msg.Data))
				var layer *natsS.LayerUpdate
				errJson := json.Unmarshal(msg.Data, &layer)
				fmt.Println("Next layer: ", layer.LayerID)
				if errJson != nil {
					fmt.Println("Error parsing json layer: ", errJson)
					msg.Nak()
					continue
				}
				saveErr := s.WriteDB.SaveLayer(layer)
				if saveErr != nil {
					fmt.Println("Failed to save layer")
					msg.Nak()
				} else {
					fmt.Println("Layer saved")
					msg.AckSync()
				}
			}
		}
	}()
}

func (s *Sink) StartAtxSink() {
	fmt.Println("Start atx sink")
	go func() {
		for {
			msgs, err := s.atxSub.Fetch(100, nats.MaxWait(360*time.Hour))
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}

			var wg sync.WaitGroup
			for _, msg := range msgs {
				go s.processAtxMessage(msg, &wg)
			}
			wg.Wait()
		}
	}()
}

func (s *Sink) processAtxMessage(msg *nats.Msg, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Atx: ", string(msg.Data))
	var atx *natsS.Atx
	errJson := json.Unmarshal(msg.Data, &atx)
	fmt.Println("Next atx: ", atx.NodeID)
	if errJson != nil {
		fmt.Println("Error parsing json atx: ", errJson)
		msg.Nak()
		return
	}
	saveErr := s.WriteDB.SaveAtx(atx)
	if saveErr != nil {
		fmt.Println("Failed to save atx")
		msg.Nak()
	} else {
		fmt.Println("Atx saved")
		msg.AckSync()
	}
}

func (s *Sink) StartTransactionResultSink() {
	fmt.Println("Start transaction result sink")

	go func() {
		for {

			msgs, err := s.transactionsResultSub.Fetch(100, nats.MaxWait(2*time.Hour))
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}
			for _, msg := range msgs {

				fmt.Println("Transaction: ", string(msg.Data))
				var transaction *natsS.Transaction
				errJson := json.Unmarshal(msg.Data, &transaction)
				fmt.Println("Next transaction: ", transaction)
				if errJson != nil {
					fmt.Println("Error parsing json transaction: ", errJson)
					msg.Nak()
					continue
				}
				saveErr := s.WriteDB.SaveTransactions(transaction, true)
				if saveErr != nil {
					fmt.Println("Failed to save transaction")
					msg.Nak()
				} else {
					fmt.Println("Transaction saved")
					msg.AckSync()
				}
			}

		}
	}()
}

func (s *Sink) StartTransactionCreatedSink() {
	fmt.Println("Start transaction created sink")

	go func() {
		for {

			msgs, err := s.transactionsCreatedSub.Fetch(100, nats.MaxWait(2*time.Hour))
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}
			for _, msg := range msgs {

				fmt.Println("Transaction: ", string(msg.Data))
				var transaction *natsS.Transaction
				errJson := json.Unmarshal(msg.Data, &transaction)
				fmt.Println("Next transaction: ", transaction)
				if errJson != nil {
					fmt.Println("Error parsing json transaction: ", errJson)
					msg.Nak()
					continue
				}
				saveErr := s.WriteDB.SaveTransactions(transaction, false)
				if saveErr != nil {
					fmt.Println("Failed to save transaction")
					msg.Nak()
				} else {
					fmt.Println("Transaction saved")
					msg.AckSync()
				}
			}

		}
	}()
}

func (s *Sink) StartMalfeasanceSink() {
	fmt.Println("Start malfeasance created sink")

	go func() {
		for {

			msgs, err := s.malfeasanceSub.Fetch(100, nats.MaxWait(8736*time.Hour))
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				continue
			}
			for _, msg := range msgs {

				fmt.Println("Malfeasance: ", string(msg.Data))
				var malfeasance *natsS.Malfeasance
				errJson := json.Unmarshal(msg.Data, &malfeasance)
				fmt.Println("Next Malfeasance: ", malfeasance)
				if errJson != nil {
					fmt.Println("Error parsing json malfeasance: ", errJson)
					msg.Nak()
					continue
				}
				saveErr := s.WriteDB.SaveMalfeasance(malfeasance)
				if saveErr != nil {
					fmt.Println("Failed to save malfeasance")
					msg.Nak()
				} else {
					fmt.Println("Malfeasance saved")
					msg.AckSync()
				}
			}

		}
	}()
}
