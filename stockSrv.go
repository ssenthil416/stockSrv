package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	HealthEndPoint    = "/health"
	StockEndPoint     = "/stock/"
	DefSymbol         = "AMEX"
	StockAPIURL       = "https://www.worldtradingdata.com/api/v1/stock?symbol="
	AddAPIToken       = "&api_token="
	SampleAPIToken    = "demo"
	HTTPClientTimeout = 15
)

type data struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	Price           string `json:"price"`
	Close_yesterday string `json:"close_yesterday"`
	Currency        string `json:"currency"`
	Market_cap      string `json:"market_cap"`
	Volume          string `json:"volume"`
	Timezone        string `json:"timezone"`
	Timezone_name   string `json:"timezone_name"`
	Gmt_offset      string `json:"gmt_offset"`
	Last_trade_time string `json:"last_trade_time"`
}

type message struct {
	Message           string `json:"message"`
	Symbols_requested int    `json:"symbols_requested"`
	Symbols_returned  int    `json:"symbols_returned"`
	Data              []data `json:"data"`
}

var (
	token string
)

func main() {
	flag.StringVar(&token, "token", "", "token : a string var")
	flag.Parse()

	//Http router
	http.HandleFunc(HealthEndPoint, GetHealthCheck)
	http.HandleFunc(StockEndPoint, GetStockDetails)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatalf("Could not start stock server: %s\n", err.Error())
	}
}

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := DefSymbol
	arrSymbol := strings.Split(r.URL.Path, StockEndPoint)

	if len(arrSymbol) == 2 && arrSymbol[1] != "" {
		symbol = arrSymbol[1]
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}

	//Get stock details
	msg, err := callStockAPI(symbol)
	if err != nil || msg.Symbols_returned == 0 {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg.Message))
		return
	}

	//Result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	md, err := json.Marshal(msg.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}
	w.Write(md)
}

func callStockAPI(symbol string) (msg message, err error) {
	var url string
	msg = message{}

	if token != "" {
		url = StockAPIURL + symbol + AddAPIToken + token
	} else {
		url = StockAPIURL + symbol + AddAPIToken + SampleAPIToken
	}

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return msg, err
	}

	// Send req using http Client
	client := &http.Client{Timeout: HTTPClientTimeout * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return msg, err
	}
	defer resp.Body.Close()

	//Read all body of http request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return msg, err
	}

	//Responce data
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
