package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type PriceQuote struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getQuoteHandler)
	http.ListenAndServe(":8080", mux)
}

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	quote, e := getPriceQuote()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOk)
	json.NewEncoder(w).Encode(quote)
}

func getPriceQuote() (*PriceQuote, error) {
	resp, e := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	var quote PriceQuote
	e = json.Unmarshal(body, &quote)
	if e != nil {
		return nil, e
	}

	return &quote, nil
}