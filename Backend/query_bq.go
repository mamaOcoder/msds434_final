package main

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// queryBasic demonstrates issuing a query and reading results.
func queryBasic(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
			return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	q := client.Query(
			"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, " +
				"(SELECT * FROM `recidivism.test`)) ")
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
			return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
			return err
	}
	if err := status.Err(); err != nil {
			return err
	}
	it, err := job.Read(ctx)
	for {
			var row []bigquery.Value
			err := it.Next(&row)
			if err == iterator.Done {
					break
			}
			if err != nil {
					return err
			}
			fmt.Fprintln(w, row)
	}
	return nil
}
