package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go-exercise/model"
	"go-exercise/view"
)

type PairRequest struct {
	Pairs []string `json:"pairs"`
}

type LTPController struct {
	DB *sql.DB
}

func (c *LTPController) HandleLTP(w http.ResponseWriter, r *http.Request) {
	var request PairRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		view.RenderError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := []view.PairResponse{}
	for _, pair := range request.Pairs {
		amount, err := model.GetLTP(c.DB, pair)
		if err != nil {
			if err == sql.ErrNoRows {
				continue // Skip pairs not found in the database
			}
			view.RenderError(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
			return
		}
		response = append(response, view.PairResponse{Pair: pair, Amount: amount})
	}

	view.RenderJSON(w, view.LTPResponse{LTP: response})
}
