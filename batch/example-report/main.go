package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/batch"
	"github.com/remiges-tech/alya/batch/pg/batchsqlc"
)

func main() {
	dbHost := "localhost"
	dbPort := 5432
	dbUser := "alyatest"
	dbPassword := "alyatest"
	dbName := "alyatest"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("error connecting db")
	}
	defer pool.Close()

	queries := batchsqlc.New(pool)

	// Initialize SlowQuery
	slowQuery := batch.SlowQuery{
		Db:      pool,
		Queries: queries,
	}
	fmt.Println(slowQuery.Queries) // just to make compiler happy while I'm developing slowquery module

	// Register the SlowQueryProcessor for the long-running report
	err = slowQuery.RegisterProcessor("LongRunningReportApp", "generateReport", &ReportProcessor{})
	if err != nil {
		fmt.Println("Failed to register SlowQueryProcessor:", err)
		return
	}

	// Submit a slow query request
	context := batch.JSONstr(`{"userId": 123, "reportType": "sales"}`)
	input := batch.JSONstr(`{"startDate": "2023-01-01", "endDate": "2023-12-31"}`)
	reqID, err := slowQuery.Submit("LongRunningReportApp", "generateReport", context, input)
	if err != nil {
		fmt.Println("Failed to submit slow query:", err)
		return
	}

	fmt.Println("Slow query submitted. Request ID:", reqID)

	// Poll for the slow query result
	for {
		status, result, messages, err := slowQuery.Done(reqID)
		if err != nil {
			fmt.Println("Error while polling for slow query result:", err)
			return
		}

		if status == batch.BatchTryLater {
			fmt.Println("Report generation in progress. Trying again in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		if status == batch.BatchSuccess {
			fmt.Println("Report generated successfully:")
			fmt.Println("Result:", result)
			break
		}

		if status == batch.BatchFailed {
			fmt.Println("Report generation failed:")
			fmt.Println("Error messages:", messages)
			break
		}
	}
}

// ReportProcessor implements the SlowQueryProcessor interface
type ReportProcessor struct{}

func (p *ReportProcessor) DoSlowQuery(initBlock any, context batch.JSONstr, input batch.JSONstr) (status batch.BatchStatus_t, result batch.JSONstr, messages []batch.ErrorMessage, outputFiles map[string]string, err error) {
	// Parse the context and input JSON
	var contextData struct {
		UserID     int    `json:"userId"`
		ReportType string `json:"reportType"`
	}
	var inputData struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}

	err = json.Unmarshal([]byte(context), &contextData)
	if err != nil {
		return batch.BatchFailed, "", nil, nil, err
	}

	err = json.Unmarshal([]byte(input), &inputData)
	if err != nil {
		return batch.BatchFailed, "", nil, nil, err
	}

	// Perform the long-running report generation logic here
	// Use the contextData and inputData to generate the report
	// Simulate a long-running operation
	time.Sleep(10 * time.Second)

	// Example output
	reportResult := fmt.Sprintf("Report generated for user %d, type %s, from %s to %s",
		contextData.UserID, contextData.ReportType, inputData.StartDate, inputData.EndDate)

	return batch.BatchSuccess, batch.JSONstr(reportResult), nil, nil, nil
}
