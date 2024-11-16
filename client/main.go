package main

import (
	"context"
	"time"
	"net/http"
	"github.com/sk8sta13/client-server-api/entity"
	"fmt"
	"errors"
	"log"
	"io/ioutil"
	"encoding/json"
	"os"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
	defer cancel()

	req, e := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if e != nil {
		panic(e)
	}

	res, e := http.DefaultClient.Do(req)
	if e != nil {
		if errors.Is(e, context.DeadlineExceeded) {
			log.Println("A request para \"http://localhost:8080/cotacao\" excedeu o tempo limite de 300 ms.")
		}

		panic(e)
	}
	defer res.Body.Close()

	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		panic(e)
	}

	var quote entity.PriceQuote
	e = json.Unmarshal(body, &quote)
	if e != nil {
		panic(e)
	}

	fmt.Println("A cotação é: %s", quote.USDBRL.Bid)

	file, e := os.Create("./client/cotacao.txt")
	if e != nil {
		log.Println("Erro na criação do arquivo \"./client/cotacao.txt\"")
	}
	defer file.Close()

	_, e = file.WriteString(fmt.Sprintf("Dólar: %s", quote.USDBRL.Bid))
}