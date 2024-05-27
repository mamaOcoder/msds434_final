package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
)

type tableSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Mode string `json:"mode,omitempty"`
}

// Function to chunk data into smaller pieces to avoid error from BQ ingest
func chunkData(data []recidData, chunkSize int) [][]recidData {
	var chunks [][]recidData
	for chunkSize < len(data) {
		data, chunks = data[chunkSize:], append(chunks, data[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, data)
	return chunks
}

func writeToBQ(trainSet []recidData, testSet []recidData) error {
	projectID := "msds434-finalproj"
	datasetID := "recidivism"

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	schema := "schema.json"

	schemaData, err := os.ReadFile(schema)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	var schemaFields []tableSchema
	if err := json.Unmarshal(schemaData, &schemaFields); err != nil {
		log.Fatalf("Failed to unmarshal schema JSON: %v", err)
	}

	// Convert the schema fields to BigQuery table schema
	var tableSchema bigquery.Schema
	for _, field := range schemaFields {
		fieldSchema := &bigquery.FieldSchema{
			Name: field.Name,
			Type: bigquery.FieldType(field.Type),
		}
		// Set the mode (REQUIRED, NULLABLE, REPEATED)
		if field.Mode == "REQUIRED" {
			fieldSchema.Required = true
		} else if field.Mode == "REPEATED" {
			fieldSchema.Repeated = true
		} // BigQuery FieldSchema is NULLABLE by default if mode is not set

		tableSchema = append(tableSchema, fieldSchema)
	}

	// Write train table
	trainTableID := "train_recid"
	train_table := client.Dataset(datasetID).Table(trainTableID)
	if err := train_table.Create(ctx, &bigquery.TableMetadata{
		Schema: tableSchema,
	}); err != nil {
		log.Fatalf("Failed to create train table: %v", err)
	}

	// Insert data into train table
	u_train := client.Dataset(datasetID).Table(trainTableID).Uploader()
	// Chunk data and upload each chunk
	chunkSize := 500 // Adjust the chunk size as needed
	chunks := chunkData(trainSet, chunkSize)
	for _, chunk := range chunks {
		if err := u_train.Put(ctx, chunk); err != nil {
			log.Fatalf("Failed to insert chunk into train table: %v", err)
		}
	}

	fmt.Printf("Inserted %d rows\n", len(trainSet))

	// if err := u_train.Put(ctx, trainSet); err != nil {
	// 	log.Fatalf("Failed to insert data into train table: %v", err)
	// }

	log.Println("Data successfully inserted into train table")

	// Write test table
	testTableID := "test_recid"
	test_table := client.Dataset(datasetID).Table(testTableID)
	if err := test_table.Create(ctx, &bigquery.TableMetadata{
		Schema: tableSchema,
	}); err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Insert data into train table
	u_test := client.Dataset(datasetID).Table(testTableID).Uploader()
	if err := u_test.Put(ctx, testSet); err != nil {
		log.Fatalf("Failed to insert data into test table: %v", err)
	}

	log.Println("Data successfully inserted into train table")

	return nil
}
