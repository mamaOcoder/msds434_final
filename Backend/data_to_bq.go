package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
)

type tableSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Mode string `json:"mode,omitempty"`
}

// Function to chunk data into smaller pieces to avoid error from BQ ingest
func chunkData(data []processedRecidData, chunkSize int) [][]processedRecidData {
	var chunks [][]processedRecidData
	for chunkSize < len(data) {
		data, chunks = data[chunkSize:], append(chunks, data[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, data)
	return chunks
}

func createOrCheckTable(ctx context.Context, client *bigquery.Client, datasetID, tableID string) error {
	// Check if the table exists
	_, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		// If the table doesn't exist, create it
		if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == 404 {
			fmt.Println("Table doesn't exist. Creating it.")

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
				switch field.Mode {
				case "REQUIRED":
					fieldSchema.Required = true
				case "REPEATED":
					fieldSchema.Repeated = true
				} // BigQuery FieldSchema is NULLABLE by default if mode is not set

				tableSchema = append(tableSchema, fieldSchema)
			}

			table := client.Dataset(datasetID).Table(tableID)
			if err := table.Create(ctx, &bigquery.TableMetadata{
				Schema: tableSchema,
			}); err != nil {
				log.Fatalf("Failed to create table: %v", err)
			}

			log.Printf("Table %s:%s has been created.", datasetID, tableID)
			return nil
		}
		return err
	}

	log.Printf("Table %s:%s already exists.", datasetID, tableID)
	return nil
}

func writeToBQ(trainSet []processedRecidData, testSet []processedRecidData) error {
	projectID := "msds434-finalproj"
	datasetID := "recidivism"

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Write train table
	trainTableID := "train_recid"
	err = createOrCheckTable(ctx, client, datasetID, trainTableID)
	if err != nil {
		log.Fatalf("Failed to create or check table: %v", err)
	}

	time.Sleep(10 * time.Second)

	// Insert data into train table
	up_train := client.Dataset(datasetID).Table(trainTableID).Uploader()

	// Chunk data and upload each chunk
	chunkSize := 500 // Adjust the chunk size as needed
	chunks := chunkData(trainSet, chunkSize)
	for _, chunk := range chunks {
		if err := up_train.Put(ctx, chunk); err != nil {
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

	err = createOrCheckTable(ctx, client, datasetID, testTableID)
	if err != nil {
		log.Fatalf("Failed to create or check table: %v", err)
	}

	time.Sleep(5 * time.Second)

	// Insert data into test table
	up_test := client.Dataset(datasetID).Table(testTableID).Uploader()
	if err := up_test.Put(ctx, testSet); err != nil {
		log.Fatalf("Failed to insert data into test table: %v", err)
	}

	log.Println("Data successfully inserted into train table")

	return nil
}
