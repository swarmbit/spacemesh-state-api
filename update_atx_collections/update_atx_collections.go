package main

import (
    "context"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://spacemesh:<password>@spacemesh-mongodb-svc.spacemesh.svc.cluster.local:27017/admin?replicaSet=spacemesh-mongodb&authMechanism=SCRAM-SHA-256")
    //clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/admin")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    
    fmt.Println("Process nodes count")
    err = processNodeCount(client)
    if err != nil {
        panic("Failed to process nodes count")
    }
    
    epochs, err := getDocumentIDs(client, "spacemesh", "atxsEpochs")

    for _, epoch := range epochs {
        fmt.Println("Process epoch atx count: ", epoch)
        
        err := processAtxCount(client, epoch)
        if err != nil {
            panic("Failed to process atx count")
        }
        
        fmt.Println("Process epoch account atx: ", epoch)
        err = processAccountAtxEpoch(client, epoch)
        if err != nil {
            panic("Failed to process account atx")
        }
        

    }

    fmt.Println("Connection to MongoDB closed.")
}

func processNodeCount(client *mongo.Client) error {
    nodesColl := client.Database("spacemesh").Collection("nodes")
    nodesCountColl := client.Database("spacemesh").Collection("nodesCount")

	nodesResult, err := nodesColl.CountDocuments(
		context.TODO(),
		bson.D{},
	)
	if err != nil {
		return err
	}
    
    _, err = nodesCountColl.UpdateOne(
        context.TODO(),
        bson.D{{Key: "_id", Value: "nodesCount"}},
        bson.D{{Key: "$set", Value: bson.D{
            {Key: "count", Value: nodesResult},
        }}},
        options.Update().SetUpsert(true),
    )
    if err != nil {
        return err
    }
    
    return nil
}
func processAtxCount(client *mongo.Client, epoch int32) error {
    atxColl := client.Database("spacemesh").Collection("atxs")
    atxsEpochsColl := client.Database("spacemesh").Collection("atxsEpochs")

	atxResult, err := atxColl.CountDocuments(
		context.TODO(),
		bson.D{
			{Key: "publishepoch", Value: epoch},
		},
	)
	if err != nil {
		return err
	}
    
    _, err = atxsEpochsColl.UpdateOne(
        context.TODO(),
        bson.D{{Key: "_id", Value: epoch}},
        bson.D{{Key: "$set", Value: bson.D{
            {Key: "totalAtx", Value: atxResult},
        }}},
        options.Update().SetUpsert(true),
    )
    if err != nil {
        return err
    }
    
    return nil
}

func processAccountAtxEpoch(client *mongo.Client, epoch int32) error {
    atxColl := client.Database("spacemesh").Collection("atxs")
    accountAtxEpochsColl := client.Database("spacemesh").Collection("accountAtxsEpochs")

    pipeline := mongo.Pipeline{
        bson.D{
            {Key: "$match", Value: bson.D{
                {Key: "publishepoch", Value: epoch},
            },
            }},
        bson.D{
            {"$group", bson.D{
                {"_id", bson.D{{"coinbase", "$coinbase"}}},
                {"totalEffectiveNumUnits", bson.D{{"$sum", "$effective_num_units"}}},
                {"totalWeight", bson.D{{"$sum", "$weight"}}},
                {"totalAtx", bson.D{{"$sum", 1}}},
            }},
        },
    }

    cursor, err := atxColl.Aggregate(
        context.TODO(),
        pipeline,
    )
	
    if err != nil {
        return err
    }
    defer cursor.Close(context.TODO())

    // Iterate through the results and append the IDs to the slice
    for cursor.Next(context.TODO()) {
        var result bson.M
        err := cursor.Decode(&result)
        if err != nil {
            return err
        }

		fmt.Println("Process account atx result: ", result)
        _, err = accountAtxEpochsColl.UpdateOne(
            context.TODO(),
            bson.D{{Key: "_id", Value: bson.M{
                "coinbase":      result["_id"].(bson.M)["coinbase"],
                "publish_epoch": epoch,
            }}},
            bson.D{{Key: "$set", Value: bson.D{
                {Key: "totalEffectiveNumUnits", Value: result["totalEffectiveNumUnits"]},
                {Key: "totalWeight", Value: result["totalWeight"]},
                {Key: "totalAtx", Value: result["totalAtx"]},
            }}},
            options.Update().SetUpsert(true),
        )
        if err != nil {
            return err
        }
    }

    if err := cursor.Err(); err != nil {
        return err
    }

    return nil
}
func getDocumentIDs(client *mongo.Client, databaseName, collectionName string) ([]int32, error) {
    // Select the database and collection
    collection := client.Database(databaseName).Collection(collectionName)

    // Create a projection to only return the _id field
    projection := bson.M{"_id": 1}

    // Execute the query
    cursor, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetProjection(projection))
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    // Create a slice to store the IDs
    var ids []int32

    // Iterate through the results and append the IDs to the slice
    for cursor.Next(context.TODO()) {
        var result bson.M
        err := cursor.Decode(&result)
        if err != nil {
            return nil, err
        }
        if id, ok := result["_id"].(int32); ok {
            ids = append(ids, id)
        }
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return ids, nil
}
