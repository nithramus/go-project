package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type scoot struct {
	BatteryLevel       string      `json:"battery_level"`
	BatteryPercent     float64     `json:"battery_percent"`
	ChargeBountyPrice  float64     `json:"charge_bounty_price"`
	Company            string      `json:"company"`
	FetchedAt          string      `json:"fetched_at"`
	LastGpsAt          string      `json:"last_gps_at"`
	LastTripAt         string      `json:"last_trip_at"`
	Latitude           float64     `json:"latitude"`
	LikelihoodScore    float64     `json:"likelihood_score"`
	Longitude          float64     `json:"longitude"`
	Mode               string      `json:"mode"`
	TaskType           string      `json:"task_type"`
	VehicleID          intOrString `json:"vehicle_id"`
	VehicleIDToDisplay intOrString `json:"vehicle_id_to_display"`
	VehicleLocationID  string      `json:"vehicle_location_id"`
	VehicleType        string      `json:"vehicle_type"`
}

type intOrString string

func (w *intOrString) UnmarshalJSON(data []byte) error {
	*w = intOrString(string(data))
	return nil
}

type coordonate struct {
	top   float64
	bot   float64
	left  float64
	right float64
}

type vehicles struct {
	Vehicles []scoot `json:"vehicles"`
}

type operator struct {
	name               string
	maxScootbyRequests int
	maxDivision        float64
}

// 48.85882,2.33068&northeast_point=48.91435,2.41599&southwest_point=48.80322,2.2453
// 48.85945,2.29862&northeast_point=48.91498,2.38389&southwest_point=48.80385,2.21335
func createURL(actualOperator operator, position coordonate) *url.URL {
	q := url.Values{}
	southwestPoint := strconv.FormatFloat(position.bot, 'f', -1, 32) + "," + strconv.FormatFloat(position.left, 'f', -1, 32)
	northeastPoint := strconv.FormatFloat(position.top, 'f', -1, 32) + "," + strconv.FormatFloat(position.right, 'f', -1, 32)
	userLocation := strconv.FormatFloat((position.top+position.bot)/2, 'f', -1, 32) + "," + strconv.FormatFloat((position.left+position.right)/2, 'f', -1, 32)
	q.Add("user_location", userLocation)
	q.Add("northeast_point", northeastPoint)
	q.Add("southwest_point", southwestPoint)
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
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", `application/json, text/plain, */*`)
	req.Header.Add("Referer", `https://scootermap.com/map`)
	req.Header.Add("Origin", `https://scootermap.com`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`)
	req.Header.Add("Sec-Fetch-Mode", `cors' --compressed &`)

	return req
}

func getScootInCoordonates(actualOperator operator, position coordonate, ch chan []scoot, client *http.Client) {
	var dat vehicles
	var req *http.Request
	req = createRequest(actualOperator, position)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		ch <- []scoot{}
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		ch <- []scoot{}
		return
	}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println(err)
		ch <- []scoot{}
		return
	}
	ch <- dat.Vehicles
	return
}

func getTrott(actualOperator operator, position coordonate, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := &http.Transport{
		IdleConnTimeout: 60 * time.Second,
	}
	client := &http.Client{Transport: tr}

	var counter int
	for true {
		counter++
		if counter == 10 {
			counter = 0
		}
		ch := make(chan []scoot)
		incrementY := (position.top - position.bot) / actualOperator.maxDivision
		incrementX := (position.right - position.left) / actualOperator.maxDivision
		var i float64
		channelCount := int(actualOperator.maxDivision * actualOperator.maxDivision)
		// we launch the process that will handle all the datas
		go insertData(ch, actualOperator, channelCount, counter)
		for i = 0; i < actualOperator.maxDivision; i++ {
			var j float64
			for j = 0; j < actualOperator.maxDivision; j++ {
				newPos := coordonate{
					position.bot + i*incrementY,
					position.bot + (i+1)*incrementY,
					position.left + j*incrementX,
					position.left + (j+1)*incrementX,
				}
				go getScootInCoordonates(actualOperator, newPos, ch, client)
			}
			time.Sleep(2 * time.Second)
		}
		time.Sleep(300 * time.Second)
	}
}

func main() {
	trottList := []operator{
		operator{"lime", 50, 60},
		operator{"bird", 50, 20},
		operator{"hive", 50, 10},
		operator{"circ", 50, 10},
		operator{"tier", 50, 10},
		operator{"voi", 50, 10},
		operator{"wind", 50, 10},
	}

	coordonateList := coordonate{48.9, 48.8, 2.20, 2.44}
	var wg sync.WaitGroup
	wg.Add(len(trottList))
	for _, operator := range trottList {
		go getTrott(operator, coordonateList, &wg)
	}
	wg.Wait()
}
