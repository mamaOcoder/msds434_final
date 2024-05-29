package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Enter ID to predict.\n"
	fmt.Fprintf(w, "%s", Body)
}

// Add function to print the metrics from the BigQuery ML model
// func getModelMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Here you can implement code to fetch and return the model metrics.

// }

// GetPredictionHandler is used to get the prediction result for a given ID
func getPredictionHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")

	recidID := paramStr[len(paramStr)-1]

	//fmt.Fprintf(w, "Prediction results for %s:\n", recidID)
	predictions, err := predictQuery(recidID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching prediction results: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing HTML template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, predictions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
