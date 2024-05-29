package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Add function to print the metrics from the BigQuery ML model
// func getModelMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Here you can implement code to fetch and return the model metrics.

// }

// GetPredictionHandler is used to get the prediction result for a given ID
func getPredictionHandler(w http.ResponseWriter, r *http.Request) {
	// paramStr := strings.Split(r.URL.Path, "/")
	// recidID := paramStr[len(paramStr)-1]
	//fmt.Fprintf(w, "Prediction results for %s:\n", recidID)

	recidID := r.URL.Query().Get("recidID")
	if recidID == "" {
		http.Error(w, "Missing recidivism ID parameter", http.StatusBadRequest)
		return
	}
	predictions, err := predictQuery(recidID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching prediction results: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("results.html")
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
