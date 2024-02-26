package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuoteResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	//Solicitando a cotação de dólar para o server
	quoteResponse, err := getQuote(ctx)
	if err != nil {
		select {
		case <-ctx.Done():
			log.Fatal("Erro: Timeout ao aguardar resposta do servidor.")
		default:
			log.Fatal("Erro: ", err)
		}
		return
	}

	//Salvando no arquivo texto
	fileName := "cotacao.txt"
	content := fmt.Sprintf("Dólar: %s\n", quoteResponse.Bid)
	err = writeFile(fileName, content)
	if err != nil {
		log.Fatal("Erro: Falha ao escrever no arquivo.")
	}

	fmt.Printf("Conteúdo recuperado e adicionado no arquivo")
}

func getQuote(ctx context.Context) (*QuoteResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
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

	var quote QuoteResponse
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil

}

func writeFile(fileName string, content string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
