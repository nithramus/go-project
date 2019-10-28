package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type mongoConnection struct {
	ctx      context.Context
	database *mongo.Database
}

func getMongoConnection() mongoConnection {
	host := os.Getenv("MONGO_HOST")
	login := os.Getenv("MONGO_LOGIN")
	password := os.Getenv("MONGO_PASSWORD")
	bdd := os.Getenv("MONGO_BDD_NAME")
	// fmt.Println(host, login, password, bdd)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + login + ":" + password + "@" + host + ":27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	database := client.Database(bdd)
	return mongoConnection{ctx: ctx, database: database}
}

func (dbConnection *mongoConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := dbConnection.ctx

	keys := r.URL.Query()
	// fromDate := keys.Get("fromDate")
	toDate := keys.Get("date")
	operator := keys.Get("operator")

	// now := time.Now()
	// new := now.Add(+24 * time.Hour)
	startDay, err := time.Parse(time.RFC1123, toDate)
	endDay := startDay.Add(24 * time.Hour)
	if err != nil {
		panic(err)
	}
	projectionDiffs := bson.M{
		"date":       1,
		"added.id":   1,
		"added.lt":   1,
		"added.lg":   1,
		"removed.id": 1,
		"removed.lt": 1,
		"removed.lg": 1,
	}
	resp, err := dbConnection.database.Collection(operator).Find(ctx, bson.M{
		"type": "diffs", "date": bson.M{"$gte": startDay, "$lte": endDay},
	}, options.Find().SetProjection(projectionDiffs))
	if err != nil {
		panic(err)
	}
	var diffs []bson.M
	for resp.Next(ctx) {
		var result bson.M
		err := resp.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		diffs = append(diffs, result)
	}
	projectionBase := bson.M{
		"scooter_list.lt": 1,
		"scooter_list.lg": 1,
		"scooter_list.id": 1,
		"date":            1,
	}
	resp, err = dbConnection.database.Collection(operator).Find(ctx, bson.M{
		"type": "base", "date": bson.M{"$gte": startDay, "$lte": endDay},
	}, options.Find().SetProjection(projectionBase))
	var base bson.M
	for resp.Next(ctx) {
		err := resp.Decode(&base)
		if err != nil {
			log.Fatal(err)
		}
	}
	json.NewEncoder(w).Encode(bson.M{"base": base, "listDiffs": diffs})
	defer resp.Close(ctx)
}

func main() {
	dbConnection := getMongoConnection()
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/data", dbConnection.ServeHTTP)
	http.Handle("/", fs)
	log.Println("Listening...")
	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		port = ":3000"
	}
	fmt.Println(port)
	err := http.ListenAndServe(port, nil)
	fmt.Println(err)
}
