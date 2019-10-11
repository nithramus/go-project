package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type mongoConnection struct {
	ctx      context.Context
	database *mongo.Database
	test     string
}

func getMongoConnection() mongoConnection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	database := client.Database("testing")
	return mongoConnection{ctx: ctx, database: database, test: "yolo"}
}

func (dbConnection *mongoConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := dbConnection.ctx
	projection := bson.M{
		"scooter_list.latitude":  1,
		"scooter_list.longitude": 1,
	}
	keys := r.URL.Query()
	// fromDate := keys.Get("fromDate")
	// toDate := keys.Get("toDate")
	operator := keys.Get("operator")

	// now := time.Now()
	// new := now.Add(+24 * time.Hour)
	fmt.Println(operator)

	resp, err := dbConnection.database.Collection(operator).Find(ctx, bson.M{
		// "date": bson.M{"$g	te": bson.Now().Add(-24 * time.Hour)},
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
