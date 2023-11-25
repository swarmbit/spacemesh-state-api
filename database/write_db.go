package database

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spacemeshos/go-scale"
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/nats"
	"github.com/swarmbit/spacemesh-state-api/pkg/transactionparser"
	"github.com/swarmbit/spacemesh-state-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WriteDB struct {
	client *mongo.Client
}

const database = "spacemesh"
const rewardsCollection = "rewards"
const layersCollection = "layers"
const atxsCollection = "atxs"
const atxsEpochsCollection = "atxsEpochs"
const nodesCollection = "nodes"
const networkInfoCollection = "networkInfo"
const accountsCollection = "accounts"
const transactionsCollection = "transactions"

func NewWriteDB(dbConnection string) (*WriteDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	createIndexes(client)
	return &WriteDB{
		client: client,
	}, err
}

func createIndexes(client *mongo.Client) error {
	rewardsColl := client.Database(database).Collection(rewardsCollection)
	rewardsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "coinbase", Value: 1},
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "node_id", Value: 1},
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err := rewardsColl.Indexes().CreateMany(context.TODO(), rewardsIndexes)
	if err != nil {
		log.Println(err)
		return err
	}

	transactionsColl := client.Database(database).Collection(transactionsCollection)
	transactionsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "principal_account", Value: 1},
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "receiver_account", Value: 1},
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "layer", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err = transactionsColl.Indexes().CreateMany(context.TODO(), transactionsIndexes)
	if err != nil {
		log.Println(err)
		return err
	}

	accountsColl := client.Database(database).Collection(accountsCollection)
	accountsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "balance", Value: -1},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err = accountsColl.Indexes().CreateMany(context.TODO(), accountsIndexes)
	if err != nil {
		log.Println(err)
		return err
	}

	atxColl := client.Database(database).Collection(atxsCollection)
	atxIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "publishepoch", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "node_id", Value: 1},
				{Key: "publishepoch", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "coinbase", Value: 1},
				{Key: "publishepoch", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "publishepoch", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "publishepoch", Value: 1},
				{Key: "effective_num_units", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err = atxColl.Indexes().CreateMany(context.TODO(), atxIndexes)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *WriteDB) SaveLayer(layer *nats.LayerUpdate) error {
	// only store processed layers
	if layer.Status > 0 {
		layersColl := m.client.Database(database).Collection(layersCollection)
		_, err := layersColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: layer.LayerID}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: layer.Status}}}},
			options.Update().SetUpsert(true),
		)
		return err
	}
	return nil
}

func (m *WriteDB) SaveAtx(atx *nats.Atx) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		atxsColl := m.client.Database(database).Collection(atxsCollection)
		atxsEpochsColl := m.client.Database(database).Collection(atxsEpochsCollection)
		nodesColl := m.client.Database(database).Collection(nodesCollection)
		accountsColl := m.client.Database(database).Collection(accountsCollection)
		weight := getATXWeight(atx.TickCount, uint64(atx.EffectiveNumUnits))
		atxDoc := &types.AtxDoc{
			AtxID:             atx.AtxID,
			NodeID:            atx.NodeID,
			EffectiveNumUnits: atx.EffectiveNumUnits,
			BaseTick:          atx.BaseTick,
			TickCount:         atx.TickCount,
			Sequence:          atx.Sequence,
			PublishEpoch:      atx.PublishEpoch,
			Coinbase:          atx.Coinbase,
			Received:          atx.Received,
			Weight:            weight,
		}
		updateResult, err := atxsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: atx.AtxID}},
			bson.D{{Key: "$set", Value: atxDoc}},
			options.Update().SetUpsert(true))
		if err != nil {
			return updateResult, err
		}

		// only update counts if inserted new ATX
		if updateResult.UpsertedCount == 1 {
			updateResult, err = atxsEpochsColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: atxDoc.PublishEpoch}},
				bson.D{{Key: "$inc", Value: bson.D{
					{Key: "totalEffectiveNumUnits", Value: atx.EffectiveNumUnits},
					{Key: "totalWeight", Value: weight},
				}}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return updateResult, err
			}

			updateResult, err = nodesColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: atxDoc.NodeID}},
				bson.D{{Key: "$addToSet", Value: bson.D{
					{Key: "atxs", Value: bson.D{
						{Key: "coinbase", Value: atxDoc.Coinbase},
						{Key: "effectiveNumUnits", Value: atxDoc.EffectiveNumUnits},
						{Key: "sequence", Value: atxDoc.Sequence},
						{Key: "weight", Value: atxDoc.Weight},
						{Key: "publishEpoch", Value: atxDoc.PublishEpoch},
						{Key: "received", Value: atxDoc.Received},
					}},
				}}},
				options.Update().SetUpsert(true),
			)

			updateResult, err = accountsColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: atxDoc.Coinbase}},
				bson.D{{Key: "$setOnInsert", Value: bson.D{
					{Key: "_id", Value: atxDoc.Coinbase},
				}}},
				options.Update().SetUpsert(true),
			)
			return updateResult, err
		}

		return updateResult, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Println("Atx transaction failed: %v", err)
	}

	fmt.Println("Atx transaction succeeded")

	return err

}

func (m *WriteDB) SaveMalfeasance(malfeasance *nats.Malfeasance) error {
	nodesColl := m.client.Database(database).Collection(nodesCollection)
	_, err := nodesColl.UpdateOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: malfeasance.NodeID}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "malfeasance", Value: bson.D{
				{Key: "received", Value: malfeasance.Received},
			}},
		}}},
		options.Update().SetUpsert(true),
	)
	fmt.Println("Malfeasance succeeded")
	return err
}

func (m *WriteDB) SaveTransactions(transaction *nats.Transaction, result bool) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		var transactionDoc *types.TransactionDoc
		if result {

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
				Complete:        true,
			}

			transactionsColl := m.client.Database(database).Collection(transactionsCollection)
			accountsColl := m.client.Database(database).Collection(accountsCollection)

			previousTransaction := transactionsColl.FindOneAndUpdate(
				context.TODO(),
				bson.D{{Key: "_id", Value: transaction.ID}},
				bson.D{{Key: "$set", Value: transactionDoc}},
				options.FindOneAndUpdate().SetUpsert(true))

			err = previousTransaction.Err()
			if err != nil && err != mongo.ErrNoDocuments {
				return previousTransaction, err
			}

			updateBalances := false

			if err == mongo.ErrNoDocuments {
				// if no documents found means it got the result before the created so it should update balance
				updateBalances = true
			} else {
				previousTransactionDoc := &types.TransactionDoc{}
				err := previousTransaction.Decode(previousTransactionDoc)
				if err != nil {
					return previousTransaction, err
				}
				// if not complete means it should update balance, if complete is a duplicate so don't update balances
				updateBalances = !previousTransactionDoc.Complete
			}

			// if transaction not sucessfull or addressess length less than 2 it means is an ineffective transaction
			if transaction.Header.Status != uint8(sTypes.TransactionSuccess) || len(transaction.Header.Addresses) < 2 {
				updateBalances = false
			}

			// if amount is 0 there is not point updating the balance for receiver account
			if updateBalances && transactionDoc.Amount > 0 {
				updateResult, err := accountsColl.UpdateOne(
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

			// update balance for sender account
			if updateBalances {

				fee := transactionDoc.Gas * transactionDoc.GasPrice
				valueToDeduct := (int64(transactionDoc.Amount) + int64(fee)) * -1
				updateResult, err := accountsColl.UpdateOne(
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
			}

			return previousTransaction, err
		} else {
			transactionDoc = &types.TransactionDoc{
				ID:              transaction.ID,
				PrincipaAccount: transaction.Header.Principal,
				Fee:             transaction.Header.Fee,
				Gas:             transaction.Header.Gas,
				Layer:           transaction.Header.LayerID,
				Status:          transaction.Header.Status,
				Method:          transaction.Header.Method,
				Complete:        false,
			}

			transactionsColl := m.client.Database(database).Collection(transactionsCollection)

			insertResult, err := transactionsColl.InsertOne(
				context.TODO(),
				transactionDoc,
			)

			// if already saved ignore error
			if err != nil && docExistsErr(err) {
				return insertResult, nil
			}

			return insertResult, err
		}
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Println("Transaction failed: %v", err)
	}

	fmt.Println("Transaction succeeded")

	return err

}

func (m *WriteDB) SaveReward(reward *nats.Reward) error {
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		rewardsColl := m.client.Database(database).Collection(rewardsCollection)
		accountsColl := m.client.Database(database).Collection(accountsCollection)
		networkInfoColl := m.client.Database(database).Collection(networkInfoCollection)

		rewardDoc := &types.RewardsDoc{
			Id:          reward.ID,
			Coinbase:    reward.Coinbase,
			LayerReward: int64(reward.LayerReward),
			TotalReward: int64(reward.Total),
			AtxID:       reward.AtxID,
			NodeId:      reward.NodeID,
			Layer:       int64(reward.Layer),
		}

		updateResult, err := rewardsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: rewardDoc.Id}},
			bson.D{{Key: "$set", Value: rewardDoc}},
			options.Update().SetUpsert(true))
		if err != nil {
			return updateResult, err
		}

		// only update counts if inserted new reward
		if updateResult.UpsertedCount == 1 {
			updateResult, err = accountsColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: reward.Coinbase}},
				bson.D{{Key: "$inc", Value: bson.D{
					{Key: "totalRewards", Value: reward.Total},
					{Key: "balance", Value: reward.Total},
				}}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return updateResult, err
			}

			updateResult, err = networkInfoColl.UpdateOne(
				context.TODO(),
				bson.D{{Key: "_id", Value: "info"}},
				bson.D{{Key: "$inc", Value: bson.D{
					{Key: "circulatingSupply", Value: reward.Total},
				}}},
				options.Update().SetUpsert(true),
			)
			return updateResult, err
		}
		return updateResult, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Println("Rewards transaction failed: %v", err)
	}

	fmt.Println("Rewards transaction succeeded")

	return err

}

func (m *WriteDB) CloseWrite() {
	m.client.Disconnect(context.TODO())
}

func getATXWeight(numUnits, tickCount uint64) uint64 {
	return safeMul(numUnits, tickCount)
}

func safeMul(a, b uint64) uint64 {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		panic("uint64 overflow")
	}
	return c
}

func docExistsErr(err error) bool {
	if wes, ok := err.(mongo.WriteException); ok {
		if wes.HasErrorCode(11000) {
			return true
		}
	}
	return false
}
