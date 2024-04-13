package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddressProvider interface {
	GetAddress(ctx context.Context, cep string) (*AddressResult, error)
}

type AddressResult struct {
	Address map[string]string
	Api     string
}

type BrasilAPI struct{}

func (api BrasilAPI) GetAddress(ctx context.Context, cep string) (*AddressResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	address, err := fetchAddress(ctx, "https://brasilapi.com.br/api/cep/v1/"+cep)
	if err != nil {
		return nil, err
	}
	return &AddressResult{Address: address, Api: "BrasilAPI"}, nil
}

type ViaCEP struct{}

func (api ViaCEP) GetAddress(ctx context.Context, cep string) (*AddressResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	address, err := fetchAddress(ctx, "http://viacep.com.br/ws/"+cep+"/json")
	if err != nil {
		return nil, err
	}
	return &AddressResult{Address: address, Api: "ViaCEP"}, nil
}

func fetchAddress(ctx context.Context, url string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar solicitação para API %s: %s", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição para API %s: %s", url, err)
	}
	defer resp.Body.Close()

	var address map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON da API %s: %s", url, err)
	}

	return address, nil
}

func main() {
	cep := "01153000"

	chAddress := make(chan *AddressResult, 2)
	go fetchAddressAsync(context.Background(), BrasilAPI{}, cep, chAddress)
	go fetchAddressAsync(context.Background(), ViaCEP{}, cep, chAddress)
	result := <-chAddress

	if result != nil {
		fmt.Printf("Resultado da API %s:\n", result.Api)
		printAddress(result.Address)
	} else {
		fmt.Println("Ambas as solicitações expiraram")
	}
}

func fetchAddressAsync(ctx context.Context, provider AddressProvider, cep string, chAddress chan<- *AddressResult) {
	addressResult, err := provider.GetAddress(ctx, cep)
	if err != nil {
		fmt.Println(err)
		chAddress <- nil
		return
	}

	chAddress <- addressResult
}

func printAddress(address map[string]string) {
	for key, value := range address {
		fmt.Printf("%s: %s\n", key, value)
	}
}
