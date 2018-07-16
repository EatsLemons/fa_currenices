package rest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Rest struct {
	httpServer *http.Server
}

func (rs *Rest) Run(port int) {
	log.Printf("[INFO] server started at :%d", port)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/currency", rs.getCurrenciesHandler).Methods("GET")

	// r.Use(rs.responseMiddleware)

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

}
