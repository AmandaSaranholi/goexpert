package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	UsdBrl UsdBrl `json:"USDBRL"`
}

type UsdBrl struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type QuoteResponse struct {
	Bid string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", quoteHandler)
	http.ListenAndServe(":8080", mux)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	ctxQuote, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	//Consumindo API externa para obter a cotação do dólar
	quote, err := getQuote(ctxQuote)
	if err != nil {
		select {
		case <-ctxQuote.Done():
			log.Println("Erro: Timeout ao chamar a API de cotação do dólar.")
			http.Error(w, "Timeout ao chamar a API de cotação do dólar", http.StatusRequestTimeout)
		default:
			log.Println("Erro: Falha ao obter a cotação do dolar.", err)
			http.Error(w, "Falha ao obter a cotação do dolar", http.StatusInternalServerError)
		}
		return
	}

	//Contando com o banco de dados
	db, err := dbPrepare()
	if err != nil {
		log.Println("Erro: Falha com o banco de dados.", err)
		http.Error(w, "Falha com o banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//Salvando cotação no banco de dados
	ctxDb, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err = saveQuote(ctxDb, db, quote)
	if err != nil {
		select {
		case <-ctxDb.Done():
			log.Println("Erro: Timeout ao persistir dados no banco.")
			http.Error(w, "Timeout ao persistir dados no banco", http.StatusRequestTimeout)
		default:
			log.Println("Erro: Falha ao inserir cotação no banco.", err)
			http.Error(w, "Falha ao inserir cotação no banco", http.StatusInternalServerError)
		}
		return
	}

	//Retornando cotação para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response QuoteResponse
	response.Bid = quote.UsdBrl.Bid
	json.NewEncoder(w).Encode(response)
}

func getQuote(ctx context.Context) (*Quote, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        code string,
		codein string,
		name string,
		high string,
		low string,
		varBid string,
		pctChange string,
		bid string,
		ask string,
		timestamp string,
		create_date string,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)

	if err != nil {
		return err
	}

	return nil
}

func dbPrepare() (*sql.DB, error) {
	db, err := dbConnect()
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveQuote(ctx context.Context, db *sql.DB, quote *Quote) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx,
		"INSERT INTO quotes (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		quote.UsdBrl.Code,
		quote.UsdBrl.CodeIn,
		quote.UsdBrl.Name,
		quote.UsdBrl.High,
		quote.UsdBrl.Low,
		quote.UsdBrl.VarBid,
		quote.UsdBrl.PctChange,
		quote.UsdBrl.Bid,
		quote.UsdBrl.Ask,
		quote.UsdBrl.Timestamp,
		quote.UsdBrl.CreateDate,
	)

	if err != nil {
		return err
	}

	return nil
}
