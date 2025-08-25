package main

import (
	"go-exercise/controller"
	"go-exercise/model"
	"log"
	"net/http"
)

func main() {
	// Initialize Kraken ticker client
	tickerClient := model.NewTickerClient()

	ltpController := &controller.LTPController{Ticker: tickerClient}
	http.HandleFunc("/api/v1/ltp", ltpController.HandleLTP)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}