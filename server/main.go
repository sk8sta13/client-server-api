package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
	"database/sql"
	"context"
	"errors"
	"github.com/sk8sta13/client-server-api/entity"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getQuoteHandler)
	http.ListenAndServe(":8080", mux)
}

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	quote, e := getPriceQuote()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote)
}

func getPriceQuote() (*entity.PriceQuote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
	defer cancel()

	req, e := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if e != nil {
		return nil, e
	}

	res, e := http.DefaultClient.Do(req)
	if e != nil {
		if errors.Is(e, context.DeadlineExceeded) {
			log.Println("A request para \"https://economia.awesomeapi.com.br/json/last/USD-BRL\" excedeu o tempo limite de 200 ms.")
		}

		return nil, e
	}
	defer res.Body.Close()

	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	var quote entity.PriceQuote
	e = json.Unmarshal(body, &quote)
	if e != nil {
		return nil, e
	}

	saveQuote(&quote)

	return &quote, nil
}

func saveQuote(quote *entity.PriceQuote) {
	DB, e := sql.Open("sqlite3", "./server/quotes.db")
	if e != nil {
		log.Fatal(e)
	}
	defer DB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	query := "INSERT INTO dollars (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, e = DB.ExecContext(
		ctx,
		query,
		quote.USDBRL.Code,
		quote.USDBRL.Codein,
		quote.USDBRL.Name,
		quote.USDBRL.High,
		quote.USDBRL.Low,
		quote.USDBRL.VarBid,
		quote.USDBRL.PctChange,
		quote.USDBRL.Bid,
		quote.USDBRL.Ask,
		quote.USDBRL.Timestamp,
		quote.USDBRL.CreateDate,
	)
	if e != nil {
		if errors.Is(e, context.DeadlineExceeded) {
			log.Println("A operação de inserção excedeu o tempo limite de 10 ms.")
		} else {
			log.Fatal(e)
		}
	}
}