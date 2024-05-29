package main

import (
	"fmt"
	"time"
)

func main() {

	// Collect full dataset from API
	fmt.Println("Collecting Data...")

	startTime := time.Now() // Record start time
	allRecid, err := getFullDataset()
	if err != nil {
		fmt.Printf("Error getting full dataset: %v\n", err)
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Time taken to fetch records: %s\n", elapsedTime)

	// Clean data and split data into train and test sets
	fmt.Printf("Preprocessing Data")
	cleanedRecid, err := cleanRecid(allRecid)
	if err != nil {
		fmt.Printf("Error imputing missing dataset: %v\n", err)
	}
	fmt.Println("Splitting Data...")
	trainSet, testSet, err := splitTrainTest(cleanedRecid)
	if err != nil {
		fmt.Printf("Error splitting dataset: %v\n", err)
	}

	fmt.Println("Test GA!!")
	// Write data to GCP BigQuery
	err = writeToBQ(trainSet, testSet)
}
