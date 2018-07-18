package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EatsLemons/fa_currencies/currency"
	"github.com/gorilla/mux"
)

type Rest struct {
	CurrService *currency.CurrencyService

	httpServer *http.Server
}

func (rs *Rest) Run(port int) {
	log.Printf("[INFO] server started at :%d", port)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/currency", rs.getCurrenciesHandler).Methods("GET")

	rs.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	err := rs.httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)
}

func (rs *Rest) getCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	response := rs.newResponseItem()
	from := r.URL.Query().Get("from")
	if from == "" {
		response.Errors = append(response.Errors, ErrorRs{Message: "from is empty"})
		rs.makeJSONResponse(w, response)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		response.Errors = append(response.Errors, ErrorRs{Message: "to is empty"})
		rs.makeJSONResponse(w, response)
		return
	}

	to = strings.ToUpper(to)
	from = strings.ToUpper(from)

	record, findErr := rs.CurrService.Storage.GetCurrPair(from, to)
	if findErr != nil {
		response.Errors = append(response.Errors, ErrorRs{Message: findErr.Error()})
		rs.makeJSONResponse(w, response)
		return
	}

	response.Result = &Rates{
		From: from,
		To:   to,
		Rate: record.To[to],
	}

	rs.makeJSONResponse(w, response)
	return
}

func (rs *Rest) makeJSONResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("[WARN] response marshaling fail %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (rs *Rest) newResponseItem() *CurrenciesResponse {
	result := CurrenciesResponse{
		Errors: make([]ErrorRs, 0),
	}

	return &result
}
