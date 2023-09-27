package database

import (
	"context"
	"time"

	"github.com/swarmbit/spacemesh-state-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReadDB struct {
	client *mongo.Client
}

func NewReadDB(dbConnection string) (*ReadDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	createIndexes(client)
	return &ReadDB{
		client: client,
	}, err
}

func (m *ReadDB) GetAccount(account string) (*types.AccountDoc, error) {
	accountsColl := m.client.Database(database).Collection(accountsCollection)
	accountResult := accountsColl.FindOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: account}},
	)
	accountDoc := &types.AccountDoc{}
	err := accountResult.Decode(accountDoc)
	if err != nil {
		return &types.AccountDoc{}, err
	}
	return accountDoc, nil
}

func (m *ReadDB) CountTransactions(account string) (int64, error) {
	transactionsColl := m.client.Database(database).Collection(transactionsCollection)

	filter := bson.M{
		"$or": []bson.M{
			{"principal_account": account},
			{"receiver_account": account},
		},
	}
	accountResult, err := transactionsColl.CountDocuments(
		context.TODO(),
		filter,
	)
	if err != nil {
		return 0, err
	}
	return accountResult, nil
}

func (m *ReadDB) CountRewards(account string) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)
	rewardsResult, err := rewardsColl.CountDocuments(
		context.TODO(),
		bson.D{
			{Key: "coinbase", Value: account},
		},
	)
	if err != nil {
		return 0, err
	}
	return rewardsResult, nil
}

func (m *ReadDB) CountRewardsLayers(account string, minLayer uint32, maxLayer uint32) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)
	filter := bson.M{
		"coinbase": account,
		"layer": bson.M{
			"$gte": minLayer,
			"$lt":  maxLayer,
		},
	}
	rewardsResult, err := rewardsColl.CountDocuments(
		context.TODO(),
		filter,
	)
	if err != nil {
		return 0, err
	}
	return rewardsResult, nil
}

func (m *ReadDB) SumRewardsLayers(account string, minLayer uint32, maxLayer uint32) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "coinbase", Value: account},
			{Key: "layer", Value: bson.D{
				{Key: "$gte", Value: minLayer},
				{Key: "$lt", Value: maxLayer},
			}},
		}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalSum", Value: bson.D{{Key: "$sum", Value: "$totalReward"}}},
		}},
	}

	cursor, err := rewardsColl.Aggregate(
		context.TODO(),
		mongo.Pipeline{match, group},
	)

	if err != nil {
		return 0, err
	}

	var results []*types.AggregationTotal
	if err = cursor.All(context.TODO(), &results); err != nil {
		return 0, err
	}

	var totalSum int64 = 0
	if len(results) > 0 {
		totalSum = results[0].TotalSum

	}
	return totalSum, nil
}

func (m *ReadDB) GetRewards(account string, skip int64, limit int64, sort int8) ([]*types.RewardsDoc, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"layer": sort})

	ctx := context.TODO()
	cursor, err := rewardsColl.Find(
		ctx,
		bson.D{
			{Key: "coinbase", Value: account},
		},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rewards []*types.RewardsDoc
	if err = cursor.All(ctx, &rewards); err != nil {
		return nil, err
	}
	return rewards, nil
}

func (m *ReadDB) GetTransactions(account string, skip int64, limit int64, sort int8) ([]*types.TransactionDoc, error) {
	transactionsColl := m.client.Database(database).Collection(transactionsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"layer": sort})

	ctx := context.TODO()
	filter := bson.M{
		"$or": []bson.M{
			{"principal_account": account},
			{"receiver_account": account},
		},
	}
	cursor, err := transactionsColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*types.TransactionDoc
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (m *ReadDB) CountAccounts() (int64, error) {
	accountsColl := m.client.Database(database).Collection(accountsCollection)

	ctx := context.TODO()
	filter := bson.M{}
	count, err := accountsColl.CountDocuments(
		ctx,
		filter,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *ReadDB) CountAtxEpoch(epoch uint64) (int64, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)
	atxResult, err := atxColl.CountDocuments(
		context.TODO(),
		bson.D{
			{Key: "publishepoch", Value: epoch},
		},
	)
	if err != nil {
		return 0, err
	}
	return atxResult, nil
}

func (m *ReadDB) GetAtxEpoch(epoch uint64) (*types.AtxEpochDoc, error) {
	atxEpochsColl := m.client.Database(database).Collection(atxsEpochsCollection)
	atxResult := atxEpochsColl.FindOne(
		context.TODO(),
		bson.D{
			{Key: "_id", Value: epoch},
		},
	)
	doc := &types.AtxEpochDoc{}
	atxResult.Decode(doc)
	return doc, nil
}

func (m *ReadDB) GetNetworkInfo() (*types.NetworkInfoDoc, error) {
	networkColl := m.client.Database(database).Collection(networkInfoCollection)
	infoResult := networkColl.FindOne(
		context.TODO(),
		bson.D{
			{Key: "_id", Value: "info"},
		},
	)
	doc := &types.NetworkInfoDoc{}
	err := infoResult.Decode(doc)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func (m *ReadDB) GetLastProcessedLayer() (*types.LayerDoc, error) {
	layersColl := m.client.Database(database).Collection(layersCollection)

	findOptions := options.Find()
	findOptions.SetLimit(1)
	findOptions.SetSort(bson.M{"_id": -1})

	ctx := context.TODO()
	filter := bson.M{}
	cursor, err := layersColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var layers []*types.LayerDoc
	if err = cursor.All(ctx, &layers); err != nil {
		return nil, err
	}
	return layers[0], nil
}

func (m *ReadDB) CloseRead() {
	m.client.Disconnect(context.TODO())
}
