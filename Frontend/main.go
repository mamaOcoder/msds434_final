package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create new Router
	router := mux.NewRouter()

	// route properly to respective handlers
	router.HandleFunc("/", rootHandler).Methods("GET")
	// router.HandleFunc("/{id:[0-9]+}", http.HandlerFunc(getPredictionHandler)).Methods("GET")
	// router.HandleFunc("/all", http.HandlerFunc(getPredictionHandler)).Methods("GET")
	//router.HandleFunc("/modelmetrics", getModelMetricsHandler).Methods("GET")

	router.HandleFunc("/predict", getPredictionHandler).Methods("GET")

	// Create new server and assign the router
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Staring Recidivism Prediction server on Port 8080")

	go func() {
		fmt.Println("Serving metrics API")

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe(":9100", h)
	}()

	// Start Server on defined port/host.
	server.ListenAndServe()

	fmt.Println("test frontend lint!!!")

}
