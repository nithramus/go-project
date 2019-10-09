package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
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

type operator struct {
	name               string
	maxScootbyRequests int
}

// 48.85882,2.33068&northeast_point=48.91435,2.41599&southwest_point=48.80322,2.2453
// 48.85945,2.29862&northeast_point=48.91498,2.38389&southwest_point=48.80385,2.21335
func createURL(actualOperator operator, position coordonate) *url.URL {
	q := url.Values{}
	q.Add("user_location", position.userLocation)
	q.Add("northeast_point", position.northeastPoint)
	q.Add("southwest_point", position.southwestPoint)
	q.Add("company", actualOperator.name)
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

func createRequest(actualOperator operator, position coordonate) *http.Request {
	u := createURL(actualOperator, position)
	fmt.Println(u.String())
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", `application/json, text/plain, */*`)
	req.Header.Add("Referer", `https://scootermap.com/map`)
	req.Header.Add("Origin", `https://scootermap.com`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`)
	req.Header.Add("Sec-Fetch-Mode", `cors' --compressed &`)

	return req
}

func getScootInCoordonates(actualOperator operator, position coordonate, ch chan []scoot, channelCountIncrement chan int) {
	var dat vehicles
	client := &http.Client{}
	var req *http.Request
	req = createRequest(actualOperator, position)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(dat.Vehicles))
	ch <- dat.Vehicles
	return
}

func getTrott(actualOperator operator, coordonateList []coordonate, wg sync.WaitGroup) {
	defer wg.Done()
	var counter int
	for true {
		counter++
		if counter == 10 {
			counter = 0
		}
		fmt.Println("test", coordonateList)
		ch := make(chan []scoot)
		channelCountIncrement := make(chan int)
		for _, position := range coordonateList {
			fmt.Println("in")
			go getScootInCoordonates(actualOperator, position, ch, channelCountIncrement)
		}
		go insertData(ch, actualOperator, len(coordonateList), channelCountIncrement, counter)
		time.Sleep(60 * time.Second)
	}
}

func main() {
	trottList := []operator{operator{"lime", 50}}
	coordonateList := []coordonate{coordonate{"48.85882,2.33068", "48.91435,2.41599", "48.80322,2.2453"}}
	var wg sync.WaitGroup
	wg.Add(len(trottList))
	for _, operator := range trottList {
		go getTrott(operator, coordonateList, wg)
	}
	wg.Wait()
}
