package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

func insertData(ch chan []scoot, actualOperator operator, channelCount int, channelCountIncrement chan int, counter int) bool {
	var scootList []scoot
	for i := 0; i < channelCount; i++ {
		select {
		case newCount := <-channelCountIncrement:
			channelCount += newCount
		case newScootList := <-ch:
			scootList = append(newScootList, newScootList...)

		}
	}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	collection := client.Database("testing").Collection(actualOperator.name)
	res, err := collection.InsertOne(ctx, bson.M{"date": bson.Now(), "operator": actualOperator.name, "scooter_list": scootList, "counter": counter})
	id := res.InsertedID
	fmt.Println(bson.M{"name": "pi", "value": 3.14159}, id)
	return true
}
