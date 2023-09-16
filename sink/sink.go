package sink

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	natsS "github.com/spacemeshos/go-spacemesh/nats"

	"github.com/nats-io/nats.go"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/node"
)

type Sink struct {
	DocDB      *database.DocDB
	NodeDB     *node.NodeDB
	layersSub  *nats.Subscription
	rewardsSub *nats.Subscription
	atxSub     *nats.Subscription
}

func NewSink(docDB *database.DocDB, nodeDB *node.NodeDB) *Sink {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		panic("Failed to connect to NATS")

	}
	js, _ := nc.JetStream()

	js.AddConsumer("layers", &nats.ConsumerConfig{
		Durable:        "state-api-process",
		DeliverSubject: "layers",
		DeliverGroup:   "state-api-process-layers",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("rewards", &nats.ConsumerConfig{
		Durable:        "state-api-process",
		DeliverSubject: "rewards",
		DeliverGroup:   "state-api-process-rewards",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	js.AddConsumer("atx", &nats.ConsumerConfig{
		Durable:        "state-api-process",
		DeliverSubject: "atx",
		DeliverGroup:   "state-api-process-atx",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})
	js.AddConsumer("transactions", &nats.ConsumerConfig{
		Durable:        "state-api-process",
		DeliverSubject: "transaction.results",
		DeliverGroup:   "state-api-process-transactions",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})
	js.AddConsumer("transactions", &nats.ConsumerConfig{
		Durable:        "state-api-process",
		DeliverSubject: "transaction.created",
		DeliverGroup:   "state-api-process-transactions",
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DeliverLastPolicy,
	})

	fmt.Println("Connect to nats stream")
	layersSub, _ := nc.QueueSubscribeSync("layers", "layers-queue")
	rewardsSub, _ := nc.QueueSubscribeSync("rewards", "rewards-queue")
	atxSub, _ := nc.QueueSubscribeSync("atx", "atx-queue")

	return &Sink{
		layersSub:  layersSub,
		rewardsSub: rewardsSub,
		atxSub:     atxSub,
		DocDB:      docDB,
		NodeDB:     nodeDB,
	}
}

func (s *Sink) StartRewardsSink() {
	fmt.Println("Start rewards sink")

	go func() {
		for {
			msg, err := s.rewardsSub.NextMsg(time.Hour)
			fmt.Println("New reward")
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				break
			}
			fmt.Println("Reward: ", string(msg.Data))

			var reward *natsS.Reward
			errJson := json.Unmarshal(msg.Data, &reward)
			fmt.Println("Next reward: ", reward.Layer)
			if errJson != nil {
				log.Fatal("Error parsing json reward: ", err)
				continue
			}
			s.DocDB.SaveReward(reward)
			fmt.Println("Reward saved")

			msg.Ack()
		}
	}()
}

func (s *Sink) StartLayersSink() {
	fmt.Println("Start layers sink")

	go func() {
		for {
			msg, err := s.layersSub.NextMsg(time.Hour)
			fmt.Println("New layers")
			if err == nats.ErrTimeout {
				fmt.Println("Error ", err.Error())
				break
			}
			fmt.Println("Layer: ", string(msg.Data))

			var layer *natsS.LayerUpdate
			errJson := json.Unmarshal(msg.Data, &layer)
			fmt.Println("Next layer: ", layer.LayerID)
			if errJson != nil {
				log.Fatal("Error parsing json layer: ", err)
				continue
			}
			s.DocDB.SaveLayer(layer)
			fmt.Println("Layer saved")

			msg.Ack()
		}
	}()
}

/*
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
*/
