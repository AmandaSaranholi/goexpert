package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(idx int, url string, requests int, results chan<- int, wg *sync.WaitGroup, log bool) {
	defer wg.Done()
	for i := 0; i < requests; i++ {
		if log {
			fmt.Printf("Worker %d: Request: %d\n", idx, i+1)
		}

		resp, err := http.Get(url)
		if err != nil {
			results <- http.StatusInternalServerError
			fmt.Printf("Worker %d: Response: %d  Status: %d\n", idx, i+1, http.StatusInternalServerError)
			continue
		}

		results <- resp.StatusCode
		if log {
			fmt.Printf("Worker %d: Response: %d  Status: %d\n", idx, i+1, resp.StatusCode)
		}

		resp.Body.Close()
	}
}

func gerarRelatorio(startTime time.Time, totalRequests int, statusCodes map[int]int, elapsedTime time.Duration) {
	fmt.Println("===== Relatório de Teste de Carga =====")
	fmt.Printf("Tempo total gasto na execução: %s\n", elapsedTime)
	fmt.Printf("Quantidade total de requisições realizadas: %d\n", totalRequests)
	fmt.Printf("Quantidade de requisições com status HTTP 200: %d\n", statusCodes[http.StatusOK])
	fmt.Println("Distribuição de outros códigos de status HTTP:")

	for code, count := range statusCodes {
		if code != http.StatusOK {
			fmt.Printf("Código %d: %d\n", code, count)
		}
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requisições")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	logs := flag.Bool("logs", true, "Exibir logs")

	flag.Parse()

	if *url == "" || *requests == 0 || *concurrency == 0 {
		fmt.Println("A URL do serviço deve ser fornecida.")
		flag.PrintDefaults()
		return
	}

	if *requests < *concurrency {
		fmt.Println("O total de requisições não pode ser menor que o número de chamadas simutâneas.")
		return
	}

	requestsPerWorker := *requests / *concurrency
	results := make(chan int, *requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(i, *url, requestsPerWorker, results, &wg, *logs)
	}

	wg.Wait()
	close(results)

	elapsedTime := time.Since(startTime)

	totalRequests := 0
	statusCodes := make(map[int]int)

	for status := range results {
		totalRequests++
		statusCodes[status]++
	}

	gerarRelatorio(startTime, totalRequests, statusCodes, elapsedTime)
}
