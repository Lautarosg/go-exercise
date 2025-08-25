package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"go-exercise/model"
	"go-exercise/view"
)

type PairRequest struct {
	Pairs []string `json:"pairs"`
}

type LTPController struct {
	Ticker *model.TickerClient
}

// HandleLTP handles POST requests with JSON body {"pairs":["BTC/USD","BTC/CHF"]} or empty list.
// Alternatively supports comma-separated pairs in header: X-Pairs: BTC/USD,BTC/CHF
func (c *LTPController) HandleLTP(w http.ResponseWriter, r *http.Request) {
	var pairs []string

	// Header takes precedence if present
	headerPairs := r.Header.Get("X-Pairs")
	if headerPairs != "" {
		for _, p := range strings.Split(headerPairs, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				pairs = append(pairs, p)
			}
		}
	} else {
		var request PairRequest
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&request) // ignore decode error -> treat as empty
		}
		pairs = request.Pairs
	}

	if len(pairs) == 0 { // empty list -> respond with empty ltp array
		view.RenderJSON(w, view.LTPResponse{LTP: []view.PairResponse{}})
		return
	}

	prices, err := c.Ticker.FetchLastPrices(pairs)
	if err != nil {
		view.RenderError(w, err.Error(), http.StatusBadGateway)
		return
	}

	resp := make([]view.PairResponse, 0, len(prices))
	for _, p := range pairs { // preserve requested order
		if price, ok := prices[p]; ok {
			resp = append(resp, view.PairResponse{Pair: p, Amount: price})
		}
	}
	view.RenderJSON(w, view.LTPResponse{LTP: resp})
}
