package database

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-spacemesh/nats"
	"github.com/swarmbit/spacemesh-state-api/pkg/transactionparser"
	"github.com/swarmbit/spacemesh-state-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocDB struct {
	client *mongo.Client
}

const database = "spacemesh"
const rewardsCollection = "rewards"
const layersCollection = "layers"
const epochsCollection = "epochs"
const atxsCollection = "atxs"
const accountsCollection = "accounts"
const transactionsCollection = "transactions"

func NewDocDB(dbConnection string) (*DocDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	return &DocDB{
		client: client,
	}, err
}

func (m *DocDB) SaveLayer(layer *nats.LayerUpdate) error {
	layersColl := m.client.Database(database).Collection(layersCollection)
	_, err := layersColl.UpdateOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: layer.LayerID}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: layer.Status}}}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *DocDB) SaveAtx(atx *nats.Atx) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		atxsColl := m.client.Database(database).Collection(atxsCollection)
		atxDoc := &types.AtxDoc{
			AtxID:             atx.AtxID,
			NodeID:            atx.NodeID,
			EffectiveNumUnits: atx.EffectiveNumUnits,
			BaseTick:          atx.BaseTick,
			TickCount:         atx.TickCount,
			Sequence:          atx.Sequence,
			Received:          atx.Received,
		}
		updateResult, err := atxsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: atx.AtxID}},
			bson.D{{Key: "$set", Value: atxDoc}},
			options.Update().SetUpsert(true))

		/*
			epochsColl := m.client.Database(database).Collection(epochsCollection)
			updateResult, err := atxsColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: atx.}},
				bson.D{{Key: "$set", Value: atxDoc}},
				options.Update().SetUpsert(true))
		*/
		return updateResult, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Fatalf("Atx transaction failed: %v", err)
	}

	fmt.Println("Atx transaction succeeded")

	return err

}

func (m *DocDB) SaveTransactions(transaction *nats.Transaction) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		var transactionDoc *types.TransactionDoc
		if len(transaction.Raw) > 0 {

			decoder := scale.NewDecoder(bytes.NewReader(transaction.Raw))
			parsedTransaction, err := transactionparser.Parse(decoder, transaction.Raw, uint32(transaction.Header.Method))
			if err != nil {
				fmt.Println("Failed to parse transaction: ", err)
				return nil, err
			}
			receiver := parsedTransaction.GetReceiver()
			receiverString := ""
			if len(receiver.Bytes()) > 0 {
				receiverString = receiver.String()
			}
			transactionDoc = &types.TransactionDoc{
				ID:              transaction.ID,
				PrincipaAccount: transaction.Header.Principal,
				ReceiverAccount: receiverString,
				Fee:             transaction.Header.Fee,
				Gas:             transaction.Header.Gas,
				Layer:           transaction.Header.LayerID,
				Status:          transaction.Header.Status,
				Method:          transaction.Header.Method,
				Amount:          parsedTransaction.GetAmount(),
				Counter:         parsedTransaction.GetCounter(),
				GasPrice:        parsedTransaction.GetGasPrice(),
			}
		} else {
			transactionDoc = &types.TransactionDoc{
				ID:              transaction.ID,
				PrincipaAccount: transaction.Header.Principal,
				Fee:             transaction.Header.Fee,
				Gas:             transaction.Header.Gas,
				Layer:           transaction.Header.LayerID,
				Status:          transaction.Header.Status,
				Method:          transaction.Header.Method,
			}
		}

		transactionsColl := m.client.Database(database).Collection(transactionsCollection)
		accountsColl := m.client.Database(database).Collection(accountsCollection)

		updateResult, err := transactionsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: transaction.ID}},
			bson.D{{Key: "$set", Value: transactionDoc}},
			options.Update().SetUpsert(true))
		if err != nil {
			return updateResult, err
		}

		if transactionDoc.Amount > 0 {
			updateResult, err = accountsColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: transactionDoc.ReceiverAccount}},
				bson.D{{Key: "$inc", Value: bson.D{
					{Key: "balance", Value: transactionDoc.Amount},
					{Key: "received", Value: transactionDoc.Amount},
				}}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return updateResult, err
			}
		}

		fee := transactionDoc.Gas * transactionDoc.GasPrice
		valueToDeduct := (int64(transactionDoc.Amount) + int64(fee)) * -1

		updateResult, err = accountsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: transactionDoc.PrincipaAccount}},
			bson.D{{Key: "$inc", Value: bson.D{
				{Key: "balance", Value: valueToDeduct},
				{Key: "sent", Value: transactionDoc.Amount},
				{Key: "fees", Value: fee},
			}}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return updateResult, err
		}

		return updateResult, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	fmt.Println("Transaction succeeded")

	return err

}

func (m *DocDB) SaveReward(reward *nats.Reward) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		rewardsColl := m.client.Database(database).Collection(rewardsCollection)
		accountsColl := m.client.Database(database).Collection(accountsCollection)

		rewardDoc := types.RewardsDoc{
			Coinbase:    reward.Coinbase,
			LayerReward: int64(reward.LayerReward),
			TotalReward: int64(reward.Total),
			AtxID:       reward.AtxID,
			NodeId:      reward.NodeID,
			Layer:       int64(reward.Layer),
		}

		insertResut, err := rewardsColl.InsertOne(context.TODO(), rewardDoc)
		if err != nil {
			return insertResut, err
		}

		updateResult, err := accountsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: reward.Coinbase}},
			bson.D{{Key: "$inc", Value: bson.D{
				{Key: "totalRewards", Value: reward.Total},
				{Key: "balance", Value: reward.Total},
			}}},
			options.Update().SetUpsert(true),
		)
		return updateResult, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Fatalf("Rewards transaction failed: %v", err)
	}

	fmt.Println("Rewards transaction succeeded")

	return err

}

func (m *DocDB) Close() {
	m.client.Disconnect(context.TODO())
}
