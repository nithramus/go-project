package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type mongoConnection struct {
	ctx        context.Context
	collection *mongo.Collection
	test       string
}

func getMongoConnection() mongoConnection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	collection := client.Database("testing").Collection("lime")
	return mongoConnection{ctx: ctx, collection: collection, test: "yolo"}
}

func (dbConnection *mongoConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := dbConnection.ctx
	projection := bson.M{
		"scooter_list.latitude":  1,
		"scooter_list.longitude": 1,
	}
	keys := r.URL.Query()
	from := keys.Get("from")
	now := time.Now()
	new := now.Add(+24 * time.Hour)
	fmt.Println(from, new)

	resp, err := dbConnection.collection.Find(ctx, bson.M{
		"date": bson.M{"$gte": bson.Now().Add(-24 * time.Hour)},
	}, options.Find().SetProjection(projection))
	if err != nil {
		panic(err)
	}
	var counter int
	var test []bson.M
	for resp.Next(ctx) {
		counter++
		var result bson.M
		err := resp.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		test = append(test, result)
	}
	json.NewEncoder(w).Encode(test)
	fmt.Println(counter)
	defer resp.Close(ctx)
}

func main() {
	dbConnection := getMongoConnection()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/data", dbConnection.ServeHTTP)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
