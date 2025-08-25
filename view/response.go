package view

import (
	"encoding/json"
	"net/http"
)

type PairResponse struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type LTPResponse struct {
	LTP []PairResponse `json:"ltp"`
}

func RenderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RenderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
