package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type mongoData struct {
	ScootMap map[string]scoot `json:"ScootMap"`
	date     string
}

type base struct {
	Operator string `json:"operator"`
}

type changeToMap struct {
	addedScoots  []scoot
	removedScoot []scoot
	updatedScoot []scoot
}

func getActualData(ctx context.Context, collection *mongo.Collection, op operator) mongoData {
	var scoopmap mongoData
	err := collection.FindOne(ctx, bson.M{"type": "actual"}).Decode(&scoopmap)
	fmt.Println(err)
	return scoopmap
}

func getBase(ctx context.Context, collection *mongo.Collection, formattedDate string, op operator) *base {
	var base base
	err := collection.FindOne(ctx, bson.M{"day": formattedDate, "type": "base"}).Decode(&base)
	if base.Operator == "" {
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}
	return &base
}

// func insertActiveData(ctx context.Context, op operator, collection *mongo.Collection, formattedDate string, data map[string]scoot) {
// 	collection.InsertOne(ctx, bson.M{"date": formattedDate, "type": "actual", "Scootmap": data})
// }

func upsertActiveData(ctx context.Context, collection *mongo.Collection, data map[string]scoot) {
	query := bson.M{"type": "actual"}
	update := bson.M{"$set": bson.M{"type": "actual", "ScootMap": data}}

	options := options.UpdateOptions{}
	options.SetUpsert(true)
	_, err := collection.UpdateOne(ctx, query, update, &options)
	fmt.Println(err)
}

func insertBase(ctx context.Context, collection *mongo.Collection, formattedDate string, op operator, scootList []scoot, counter int) {
	fmt.Println(op.name)
	data := bson.M{"type": "base", "day": formattedDate, "date": bson.Now(), "operator": op.name, "scooter_list": scootList, "counter": counter}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(bson.M{"name": "pi", "value": 3.14159}, id)
}

func getDiffs(oldMap map[string]scoot, newMap map[string]scoot) changeToMap {
	var listOfChanges changeToMap
	for _, scoot := range newMap {
		if _, ok := oldMap[string(scoot.Id)]; ok {
			// if oldMapScoot.Latitude != scoot.Latitude || oldMapScoot.Longitude != scoot.Longitude {
			// 	listOfChanges.updatedScoot = append(listOfChanges.updatedScoot, scoot)
			// }
		} else {
			listOfChanges.addedScoots = append(listOfChanges.addedScoots, scoot)
		}
	}
	for _, scoot := range oldMap {
		if _, ok := newMap[string(scoot.Id)]; ok {

		} else {
			listOfChanges.removedScoot = append(listOfChanges.removedScoot, scoot)
		}
	}
	return listOfChanges
}

func insertDiffs(ctx context.Context, collection *mongo.Collection, formattedDate string, op operator, listOfChanges changeToMap) {
	data := bson.M{
		"type": "diffs", "day": formattedDate,
		"date":     bson.Now(),
		"operator": op.name,
		"added":    listOfChanges.addedScoots,
		"removed":  listOfChanges.removedScoot}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}

}

func insertData(ch chan []scoot, op operator, channelCount int, counter int) bool {
	var scootMap map[string]scoot
	scootMap = make(map[string]scoot)
	for i := 0; i < channelCount; i++ {
		fmt.Println(i, channelCount)
		newScootList := <-ch
		for _, scoot := range newScootList {
			scootMap[string(scoot.Id)] = scoot
		}
	}
	fmt.Println(len(scootMap))
	var scootList []scoot
	for _, scoot := range scootMap {
		scootList = append(scootList, scoot)
	}
	host := os.Getenv("MONGO_HOST")
	login := os.Getenv("MONGO_LOGIN")
	password := os.Getenv("MONGO_PASSWORD")
	bdd := os.Getenv("MONGO_BDD_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + login + ":" + password + "@" + host + ":27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	collection := client.Database(bdd).Collection(op.name)
	day := time.Now()
	formattedDate := day.Format("02 Jan 2006")
	if err != nil {
		panic(err)
	}
	base := getBase(ctx, collection, formattedDate, op)
	if base == nil {
		fmt.Println("insert base")
		insertBase(ctx, collection, formattedDate, op, scootList, counter)
	} else {
		oldMap := getActualData(ctx, collection, op)
		fmt.Println(len(oldMap.ScootMap))
		diffs := getDiffs(oldMap.ScootMap, scootMap)
		fmt.Println("added", len(diffs.addedScoots))
		fmt.Println("removed", len(diffs.removedScoot))
		insertDiffs(ctx, collection, formattedDate, op, diffs)
	}
	upsertActiveData(ctx, collection, scootMap)
	return true
}
