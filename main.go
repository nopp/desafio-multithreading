package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type brasilApi struct {
	Cep    string `json:"cep"`
	State  string `json:"state"`
	City   string `json:"city"`
	Street string `json:"street"`
}

type viaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
}

func main() {

	argsWithProg := os.Args

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	c1 := make(chan brasilApi)
	c2 := make(chan viaCep)

	// Brasil API
	go func() {
		req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+argsWithProg[1], nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição:", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Erro ao fazer requisição:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Erro na resposta da API Brasil:", resp.Status)
			return
		}

		var data brasilApi
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Erro ao decodificar resposta da API Brasil:", err)
			return
		}

		c1 <- data

	}()

	// ViaCEP
	go func() {
		req, err := http.NewRequestWithContext(ctx, "GET", "http://viacep.com.br/ws/"+argsWithProg[1]+"/json/", nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição:", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Erro ao fazer requisição:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Erro na resposta da API ViaCEP:", resp.Status)
			return
		}

		var data viaCep
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Erro ao decodificar resposta da API ViaCEP:", err)
			return
		}

		c2 <- data
	}()

	// Selecionando a resposta mais rápida
	select {
	case msg := <-c1: // Recebendo mensagem do Brasil API
		fmt.Printf("Recebido do Brasil API: CEP: %s - %s\n", msg.Cep, msg.Street)
	case msg := <-c2: // Recebendo mensagem do ViaCEP
		fmt.Printf("Recebido do ViaCEP: CEP: %s - %s\n", msg.Cep, msg.Logradouro)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: Nenhuma resposta recebida em 1 segundo.")
	}

}
