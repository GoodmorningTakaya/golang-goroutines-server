package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

//global channel
var data = make(chan loggedData)

//server data
type lastFetched struct {
	lastFetchedTime int64
	requestsMade    uint32
}

//data sent to the global channel to be logged
type loggedData struct {
	requestTime int64
	requestIP   string
	currentTime int64
}

//instantiate server data
var recentTime lastFetched

//data storing response from worldtimeapi.org
type worldTime struct {
	Week         int    `json: "week_number"`
	UtcOffset    string `json: "utc_offset"`
	UtcDateTime  string `json: "utc_datetime"`
	UnixTime     int64  `json: "unixtime"`
	TimeZone     string `json: "timezone"`
	RawOffset    int    `json: "raw_offset"`
	DstUntil     string `json: "dst_until"`
	DstOffset    int    `json: "dst_offset"`
	DstFrom      string `json: "dst_from"`
	Dst          bool   `json: "dst"`
	DayOfYear    int    `json: "day_of_year"`
	DayOfWeek    int    `json: "day_of_week"`
	DateTime     string `json: "datetime"`
	ClientIP     string `json: "client_ip"`
	Abbreviation string `json: "abbreviation"`
}

//Makes a request to worldtimeapi.org and unmarshals the data into the server's data
func requestTime() {

	response, err := http.Get("http://worldtimeapi.org/api/ip?type:integer")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(response.Body)

		var lastretrievedtime worldTime
		jsonErr := json.Unmarshal([]byte(data), &lastretrievedtime)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		recentTime.lastFetchedTime = lastretrievedtime.UnixTime

	}
}

//Infinitely loops at the rate of Euler's constant to update the server's data
func backgroundTickTime() {
	t := float64(math.E)
	for range time.Tick(time.Duration(int64(t)) * time.Second) {
		requestTime()
		recentTime.requestsMade++

	}
}

//receives data captured at the time a GET request had been last made stored in the global channel
//writes out the data captured into a file "log.txt"
func logData(data chan loggedData) {
	logged := <-data
	logged.currentTime = time.Now().Unix()
	payload := fmt.Sprintf("%s-%s-%s\n", logged.requestIP, strconv.FormatInt(logged.currentTime, 10), strconv.FormatInt(logged.requestTime, 10))
	fo, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer fo.Close()
	if _, err := fo.WriteString(payload); err != nil {
		log.Println(err)
	}
}

//Handles the http requests made to the server, retrieving the IP making the request and the time of the request
//This information is then sent into the global channel to be logged
func root(h http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var logged loggedData
		logged.requestIP = r.Header.Get("X-Real-Ip")
		if logged.requestIP == "" {
			logged.requestIP = r.Header.Get("X-Forwarded-For")
		}
		if logged.requestIP == "" {
			logged.requestIP = r.RemoteAddr
		}
		logged.requestTime = recentTime.lastFetchedTime
		go logData(data)
		data <- logged
	}

}

func main() {
	go backgroundTickTime()
	http.HandleFunc("/", root)
	http.ListenAndServe(":12345", nil)

}
