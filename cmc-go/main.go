package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// loads postgresql drivers
	_ "github.com/lib/pq"
)

type Response struct {
	Data map[string]Currency `json:"data"`
}

type Currency struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Quotes map[string]struct {
		Price float64 `json:"price"`
	} `json:"quotes"`
}

const cmcUrl = "https://api.coinmarketcap.com/v2/ticker/?limit=10"

func getCurrencies(c *http.Client) (map[string]Currency, error) {
	resp, err := c.Get(cmcUrl)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}

	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("can't decode response: %v", err)
	}

	return r.Data, nil
}

func saveCurrencies(db *sql.DB, currencies map[string]Currency) error {
	stmt, err := db.Prepare("INSERT INTO currencies (name, symbol, price, time) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range currencies {
		_, err = stmt.Exec(c.Name, c.Symbol, c.Quotes["USD"].Price, time.Now())
	}

	return err
}

func main() {
	c := &http.Client{Timeout: 10 * time.Second}

	currencies, err := getCurrencies(c)
	if err != nil {
		log.Fatalf("can't get currencies: %v", err)
	}

	conn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBNAME"))

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS currencies (id SERIAL PRIMARY KEY,
							 name VARCHAR(64),
							 symbol VARCHAR(8),
							 price FLOAT,
							 time timestamp)`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("can't create table: %v", err)
	}

	if err := saveCurrencies(db, currencies); err != nil {
		log.Fatalf("can't save currencies: %v", err)
	}
}
