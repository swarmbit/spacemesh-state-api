package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spacemeshos/go-spacemesh/nats"
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
const accountsCollection = "accounts"
const trackingCollection = "tracking"

func NewDocDB(dbConnection string) (*DocDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	return &DocDB{
		client: client,
	}, err
}

func (m *DocDB) GetOffset(id string) (int64, error) {
	trackingColl := m.client.Database(database).Collection(trackingCollection)
	filter := bson.D{{Key: "_id", Value: id}}
	var tracking types.TrackingDoc
	err := trackingColl.FindOne(context.TODO(), filter).Decode(&tracking)
	if err == mongo.ErrNoDocuments {

		tracking := types.TrackingDoc{
			Id:     id,
			Offset: 0,
		}
		_, err := trackingColl.InsertOne(context.TODO(), tracking)
		if err != nil {
			log.Fatalf("Failed to insert document: %v", err)
			return 0, err
		}
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return tracking.Offset, nil
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

func (m *DocDB) SaveAccounts(offset int64, accounts []*types.NodeAccount) error {
	if len(accounts) > 0 {
		session, err := m.client.StartSession()
		defer session.EndSession(context.TODO())

		callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
			accountsColl := m.client.Database(database).Collection(accountsCollection)

			accountUpdatesDoc := make([]mongo.WriteModel, len(accounts))
			for i, a := range accounts {
				account := types.AccountDoc{
					Address:      a.Address.String(),
					Balance:      a.Balance,
					NextNonce:    a.NextNonce,
					LayerUpdated: int64(a.LayerUpdated),
					Template:     a.Template,
					State:        a.State,
				}
				address := a.Address.String()
				accountUpdatesDoc[i] = mongo.NewUpdateOneModel().
					SetFilter(bson.D{{Key: "_id", Value: address}}).
					SetUpdate(bson.D{{Key: "$set", Value: account}}).SetUpsert(true)

			}

			bulkWrite, err := accountsColl.BulkWrite(context.TODO(), accountUpdatesDoc)
			if err != nil {
				return bulkWrite, err
			}

			trackingColl := m.client.Database(database).Collection(trackingCollection)
			filter := bson.D{{Key: "_id", Value: "accounts"}}
			replace, err := trackingColl.ReplaceOne(context.TODO(), filter, &types.TrackingDoc{
				Id:     "accounts",
				Offset: offset + int64(len(accounts)),
			})
			return replace, err
		}

		// Execute the operations in a transaction
		if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
			log.Fatalf("Accounts transaction failed: %v", err)
		}

		fmt.Println("Accounts transaction succeeded")

		return err
	}

	fmt.Println("No accounts to add")
	return nil
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
			bson.D{{Key: "$inc", Value: bson.D{{Key: "totalRewards", Value: reward.Total}}}},
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
