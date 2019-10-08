package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type scoot struct {
	BatteryLevel       string  `json:"battery_level"`
	BatteryPercent     int     `json:"battery_percent"`
	ChargeBountyPrice  float64 `json:"charge_bounty_price"`
	Company            string  `json:"company"`
	FetchedAt          string  `json:"fetched_at"`
	LastGpsAt          string  `json:"last_gps_at"`
	LastTripAt         string  `json:"last_trip_at"`
	Latitude           float64 `json:"latitude"`
	LikelihoodScore    float64 `json:"likelihood_score"`
	Longitude          float64 `json:"longitude"`
	Mode               string  `json:"mode"`
	TaskType           string  `json:"task_type"`
	VehicleID          string  `json:"vehicle_id"`
	VehicleIDToDisplay string  `json:"vehicle_id_to_display"`
	VehicleLocationID  string  `json:"vehicle_location_id"`
	VehicleType        string  `json:"vehicle_type"`
}

type coordonate struct {
	userLocation   string
	northeastPoint string
	southwestPoint string
}

type vehicles struct {
	Vehicles []scoot `json:"vehicles"`
}

// 48.85882,2.33068&northeast_point=48.91435,2.41599&southwest_point=48.80322,2.2453
// 48.85945,2.29862&northeast_point=48.91498,2.38389&southwest_point=48.80385,2.21335
func createURL(operator string) *url.URL {
	q := url.Values{}
	q.Add("user_location", "48.85882,2.33068")
	q.Add("northeast_point", "48.91435,2.41599")
	q.Add("southwest_point", "48.80322,2.2453")
	q.Add("company", operator)
	q.Add("mode", "ride")
	q.Add("randomize", "false")

	// var u *URL
	u := &url.URL{
		Scheme:   "https",
		Host:     "vehicles.scootermap.com",
		Path:     "api/vehicles",
		RawQuery: q.Encode(),
	}
	return u
}

func createRequest(operator string) *http.Request {
	u := createURL(operator)
	fmt.Println(u.String())
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", `application/json, text/plain, */*`)
	req.Header.Add("Referer", `https://scootermap.com/map`)
	req.Header.Add("Origin", `https://scootermap.com`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`)
	req.Header.Add("Sec-Fetch-Mode", `cors' --compressed &`)

	return req
}

func getScootInCoordonates(operator string) []byte {
	client := &http.Client{}
	var req *http.Request
	req = createRequest(operator)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func insertData(old_list_scoot []scoot, list_scoot []scoot, operator string, counter int) bool {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	collection := client.Database("testing").Collection("ride")
	res, err := collection.InsertOne(ctx, bson.M{"date": bson.Now(), "operator": operator, "scooter_list": list_scoot, "counter": counter})
	id := res.InsertedID
	fmt.Println(bson.M{"name": "pi", "value": 3.14159}, id)
	return true
}

func getTrott(operator string, coordonateList []coordonate) {
	var scoots []scoot = nil
	var counter int
	for true {
		counter++
		if counter == 10 {
			counter = 0
		}
		ch := make(chan int)
		for _, position := range coordonateList {
			var dat vehicles
			body := getScootInCoordonates(operator, position)
			err := json.Unmarshal(body, &dat)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(len(dat.Vehicles))
			scoots = dat.Vehicles
			time.Sleep(60 * time.Second)
		}

		insertData(scoots, dat.Vehicles, operator, counter)
	}
}

func main() {
	trottList := []string{"bird"}
	coordonateList := []coordonate{coordonate{"48.85882,2.33068", "48.91435,2.41599", "48.80322,2.2453"}}
	for _, operator := range trottList {
		go getTrott(operator, coordonateList)
	}
}
