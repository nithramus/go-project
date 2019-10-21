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
	fmt.Println(host, login, password, bdd)
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
	projection := bson.M{
		"scooter_list.latitude":  1,
		"scooter_list.longitude": 1,
		"date":                   1,
	}
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
	fmt.Println("date:", startDay)
	fmt.Println("yolo:", endDay)
	fmt.Println(operator)

	resp, err := dbConnection.database.Collection(operator).Find(ctx, bson.M{
		"date": bson.M{"$gte": startDay, "$lte": endDay},
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
