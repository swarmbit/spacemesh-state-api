package database

import (
	"context"
	"errors"
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
		if err == mongo.ErrNoDocuments {
			return &types.AccountDoc{}, nil
		}
		return &types.AccountDoc{}, err
	}
	return accountDoc, nil
}

func (m *ReadDB) GetNode(nodeId string) (*types.NodeDoc, error) {
	nodesColl := m.client.Database(database).Collection(nodesCollection)
	nodeResult := nodesColl.FindOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: nodeId}},
	)
	nodeDoc := &types.NodeDoc{}
	err := nodeResult.Decode(nodeDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &types.NodeDoc{}, nil
		}
		return &types.NodeDoc{}, err
	}
	return nodeDoc, nil
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

func (m *ReadDB) CountRewards(account string, firstLayer int, lastLayer int) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	filter := bson.D{
		{Key: "coinbase", Value: account},
	}
	if firstLayer > -1 && lastLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$gte", firstLayer}}},
			{"layer", bson.D{{"$lte", lastLayer}}},
		}
	} else if firstLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$gte", firstLayer}}},
		}
	} else if lastLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$lte", lastLayer}}},
		}
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

func (m *ReadDB) CountNodeRewards(node string) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)
	rewardsResult, err := rewardsColl.CountDocuments(
		context.TODO(),
		bson.D{
			{Key: "node_id", Value: node},
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

func (m *ReadDB) CountNodeRewardsLayers(node string, minLayer uint32, maxLayer uint32) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)
	filter := bson.M{
		"node_id": node,
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

func (m *ReadDB) SumNodeRewardsLayers(node string, minLayer uint32, maxLayer uint32) (int64, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "node_id", Value: node},
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

func (m *ReadDB) GetRewards(account string, skip int64, limit int64, sort int8, firstLayer int, lastLayer int) ([]*types.RewardsDoc, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"layer": sort})

	filter := bson.D{
		{Key: "coinbase", Value: account},
	}
	if firstLayer > -1 && lastLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$gte", firstLayer}}},
			{"layer", bson.D{{"$lte", lastLayer}}},
		}
	} else if firstLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$gte", firstLayer}}},
		}
	} else if lastLayer > -1 {
		filter = bson.D{
			{Key: "coinbase", Value: account},
			{"layer", bson.D{{"$lte", lastLayer}}},
		}
	}

	ctx := context.TODO()
	cursor, err := rewardsColl.Find(
		ctx,
		filter,
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

func (m *ReadDB) GetNodeRewards(node string, skip int64, limit int64, sort int8) ([]*types.RewardsDoc, error) {
	rewardsColl := m.client.Database(database).Collection(rewardsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"layer": sort})

	ctx := context.TODO()
	cursor, err := rewardsColl.Find(
		ctx,
		bson.D{
			{Key: "node_id", Value: node},
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

func (m *ReadDB) GetAtxWeightAccount(account string, epoch uint64) (*types.AggregationAtxTotals, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "coinbase", Value: account},
			{Key: "publishepoch", Value: epoch},
		}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalWeight", Value: bson.D{{Key: "$sum", Value: "$weight"}}},
			{Key: "totalEffectiveNumUnits", Value: bson.D{{Key: "$sum", Value: "$effective_num_units"}}},
		}},
	}

	cursor, err := atxColl.Aggregate(
		context.TODO(),
		mongo.Pipeline{match, group},
	)

	if err != nil {
		return nil, err
	}

	var results []*types.AggregationAtxTotals
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return &types.AggregationAtxTotals{}, nil
}

func (m *ReadDB) GetAtxWeightNode(node string, epoch uint64) (*types.AggregationAtxTotals, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "node_id", Value: node},
			{Key: "publishepoch", Value: epoch},
		}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalWeight", Value: bson.D{{Key: "$sum", Value: "$weight"}}},
			{Key: "totalEffectiveNumUnits", Value: bson.D{{Key: "$sum", Value: "$effective_num_units"}}},
		}},
	}

	cursor, err := atxColl.Aggregate(
		context.TODO(),
		mongo.Pipeline{match, group},
	)

	if err != nil {
		return nil, err
	}

	var results []*types.AggregationAtxTotals
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return &types.AggregationAtxTotals{}, nil
}

func (m *ReadDB) GetTransactions(account string, skip int64, limit int64, sort int8, complete bool) ([]*types.TransactionDoc, error) {
	transactionsColl := m.client.Database(database).Collection(transactionsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"layer": sort})

	ctx := context.TODO()
	filter := bson.M{
		"$or": []bson.M{
			{"principal_account": account, "complete": complete},
			{"receiver_account": account, "complete": complete},
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

func (m *ReadDB) FilterAccountAtxNodesForEpoch(account string, epoch uint64, nodes []string) ([]string, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{"node_id", 1}})

	ctx := context.TODO()
	filter := bson.M{
		"coinbase":     account,
		"publishepoch": epoch,
		"node_id":      bson.M{"$in": nodes},
	}

	cursor, err := atxColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]string, 0)

	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		value, ok := result["node_id"]
		if !ok {
			return nil, errors.New("node_id not present in object")
		}
		results = append(results, value.(string))
	}
	return results, nil
}

func (m *ReadDB) CountAccountAtxEpoch(account string, epoch uint64) (int64, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	ctx := context.TODO()
	filter := bson.M{
		"coinbase":     account,
		"publishepoch": epoch,
	}

	count, err := atxColl.CountDocuments(
		ctx,
		filter,
	)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func (m *ReadDB) GetAccountAtxEpoch(account string, epoch uint64, skip int64, limit int64, sort int8) ([]*types.AtxDoc, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"received": sort})

	ctx := context.TODO()
	filter := bson.M{
		"coinbase":     account,
		"publishepoch": epoch,
	}

	cursor, err := atxColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var atx []*types.AtxDoc
	if err = cursor.All(ctx, &atx); err != nil {
		return nil, err
	}
	return atx, nil
}

func (m *ReadDB) GetAtxForEpoch(epoch uint64) ([]*types.AtxDoc, error) {
	atxColl := m.client.Database(database).Collection(atxsCollection)

	sortDoc := bson.D{
		{Key: "_id", Value: 1},
		{Key: "publishepoch", Value: 1},
	}

	findOptions := options.Find()
	findOptions.SetSort(sortDoc)

	ctx := context.TODO()
	filter := bson.M{
		"publishepoch": epoch,
	}
	cursor, err := atxColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var atx []*types.AtxDoc
	if err = cursor.All(ctx, &atx); err != nil {
		return nil, err
	}
	return atx, nil
}

func (m *ReadDB) GetMalfeasanceNodes() ([]*types.NodeDoc, error) {
	nodesColl := m.client.Database(database).Collection(nodesCollection)

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"publishepoch": -1})

	ctx := context.TODO()
	filter := bson.M{"malfeasance": bson.M{"$exists": true}}

	cursor, err := nodesColl.Find(
		ctx,
		filter,
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var node []*types.NodeDoc
	if err = cursor.All(ctx, &node); err != nil {
		return nil, err
	}
	return node, nil
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
	filter := bson.M{
		"status": 3,
	}
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
	if len(layers) > 0 {
		return layers[0], nil
	} else {
		return &types.LayerDoc{}, nil
	}
}

func (m *ReadDB) CloseRead() {
	m.client.Disconnect(context.TODO())
}
