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

	// Get a handle for your collection
	collection := client.Database("spacemesh").Collection("accounts")

	// Create an array of write models
	var models []mongo.WriteModel
	operations := []struct {
		ID      string
		Balance int64
	}{
		{"sm1qqqqqqylyl2l0zsmmax0wnutt4dwnrkcwef5eeq3xladz", 2743200000000000},
		{"sm1qqqqqqyp8ueuuh2dgrc2g6ps4xvueyjpky6rfaqnxdy97", 5867100000000000},
		{"sm1qqqqqqzgmt5vv4jgucas8vvrlu4daa4r29cunwqpv0trt", 1022800000000000},
		{"sm1qqqqqq80we5pmwztmqgpxu6xasapgn65r4xjczqxu39a2", 409000000000000},
		{"sm1qqqqqqy6anfdew2sdtvuuaffjy0l7ssu9r8vjsss5c442", 2045400000000000},
		{"sm1qqqqqqyw9lvmmayckrxlnf8u7850tsjdg8zz6dg956gxg", 270600000000000},
		{"sm1qqqqqq9a8g5act6ewmmmmmux8l570kr6l68htzsq94wg4", 4090900000000000},
		{"sm1qqqqqqrgqc65x5q6exujgjs970fvcakd790na3gsr3uu7", 333300000000000},
		{"sm1qqqqqqpc4ppx8s4gmdaa5tzg35s6l3v6ujg6hmqz3s4lc", 859100000000000},
		{"sm1qqqqqq8za0geafhj4avegdwhtaw9fmgjh07s55cufk695", 293300000000000},
		{"sm1qqqqqqpf6djx3axy7aag8zhyf84ljsulhfypfxgpw5y0u", 1990600000000000},
		{"sm1qqqqqq827v998nt99vupxlrfucdk0tapp2hjyygmn3kyd", 409100000000000},
		{"sm1qqqqqqpc55ghjq6sxf5k77yc8n82fkwhlj0jedcgw2zck", 4909100000000000},
		{"sm1qqqqqqxq54zvz484hhcnrghnqrjlw26twwld32slz3lxa", 191800000000000},
		{"sm1qqqqqqyf5uc2n8mutm3tuateu5efcm9awvrclmcm5mhdf", 2933540000000000},
		{"sm1qqqqqq99klpy92mwlfcft5lmz8q5sef2v2qvtucd9y55v", 2933540000000000},
		{"sm1qqqqqqyjpjgup8fz32cufcv2nlqrr3nyvge7akqt0daea", 2933540000000000},
		{"sm1qqqqqq8zukfwtggnfq4jaqpv6m8xgtg5ay2ezaqpr2w6y", 2933540000000000},
		{"sm1qqqqqqrhftrq9knsetema7dt0qfzgd5a20m9rcczk0gk5", 2933540000000000},
		{"sm1qqqqqqyfq5f522mmrzs4lczhaf30jh4pmqyfrzcg8vrpc", 3303792000000000},
		{"sm1qqqqqqx55z5795569fq5kym3gw2h6zp6ajeh46c5wtrzf", 455300000000000},
		{"sm1qqqqqqyvet26gqsxjt6w50nnp80jvajr3n25xzsdpxn65", 831250000000000},
		{"sm1qqqqqqzgqpjxdw77aw74f8mz540rykda4x2jgjgaca7z5", 184375000000000},
		{"sm1qqqqqq9s5l9tc87wspycr68dfagmzxplzdn7zlcymnkup", 15000000000000},
		{"sm1qqqqqqptx3mdg4gm67arv4ykau6nfy6w9v03x9s49wmru", 100000000000000},
		{"sm1qqqqqq9fwfymdr7qv0tfc3ppa4q8ara6qm7kwugw9gdme", 500000000000000},
		{"sm1qqqqqqy3fc8nvdetan6qjz5cju7h4c60mjyvdlqnlqpxu", 15688500000000000},
		{"sm1qqqqqqrt64knhuxu3kzq50ak04nrkk9yf2zxprshmvkcy", 88818783000000000},
	}

	for _, op := range operations {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": op.ID}).
			SetUpdate(bson.M{"$inc": bson.M{"balance": op.Balance}}).
			SetUpsert(true)
		models = append(models, model)
	}

	// Create a write option
	opts := options.BulkWrite().SetOrdered(false)

	// Bulk write
	results, err := collection.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Insert: %v\n", results.InsertedCount)
	fmt.Printf("Matched: %v\n", results.MatchedCount)
	fmt.Printf("Modified: %v\n", results.ModifiedCount)
	fmt.Printf("Upserted: %v\n", results.UpsertedCount)

	// Disconnect
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

