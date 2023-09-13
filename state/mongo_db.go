package state

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocDB struct {
	client *mongo.Client
}

const database = "spacemesh"
const rewardsCollection = "rewards"
const trackingCollection = "tracking"

func NewDocDB() (*DocDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	return &DocDB{
		client: client,
	}, err
}

func (m *DocDB) GetOffset(id string) (int64, error) {
	trackingColl := m.client.Database(database).Collection(trackingCollection)
	filter := bson.D{{Key: "_id", Value: id}}
	var tracking Tracking
	err := trackingColl.FindOne(context.TODO(), filter).Decode(&tracking)
	if err == mongo.ErrNoDocuments {

		tracking := Tracking{
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

func (m *DocDB) SaveRewards(offset int64, rewards []interface{}) error {

	// Start a session
	session, err := m.client.StartSession()
	defer session.EndSession(context.TODO())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		rewardsColl := m.client.Database(database).Collection(rewardsCollection)
		insertResut, err := rewardsColl.InsertMany(context.TODO(), rewards)
		if err != nil {
			panic(err)
		}

		trackingColl := m.client.Database(database).Collection(trackingCollection)
		filter := bson.D{{Key: "_id", Value: "rewards"}}
		trackingColl.ReplaceOne(context.TODO(), filter, &Tracking{
			Id:     "rewards",
			Offset: offset + int64(len(insertResut.InsertedIDs)),
		})
		return nil, err
	}

	// Execute the operations in a transaction
	if _, err := session.WithTransaction(context.TODO(), callback); err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	fmt.Println("Transaction succeeded")

	return err
}

func (m *DocDB) Close() {
	m.client.Disconnect(context.TODO())
}

type Tracking struct {
	Id     string `bson:"_id"`
	Offset int64  `bson:"offset"`
}

type RewardsDoc struct {
	NodeId   string `bson:"node_id"`
	Coinbase string `bson:"coinbase"`
	AtxID    string `bson:"atx_id"`
	Ammount  int64  `bson:"ammount"`
	Layer    int64  `bson:"layer"`
}
