package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamwg/distribution-bench/harness"
	"github.com/montanaflynn/stats"
)

func main() {
	resultsFilePath := "-"

	if len(os.Args) >= 2 {
		resultsFilePath = os.Args[1]
	}

	var resultsFile *os.File
	if resultsFilePath == "-" {
		resultsFile = os.Stdin
	} else {
		fd, err := os.Open(resultsFilePath)
		if err != nil {
			log.Fatalf("failed to open results file %q: %v", resultsFilePath, err)
		}
		defer fd.Close()
		resultsFile = fd
	}

	var results []harness.TestResult
	dec := json.NewDecoder(resultsFile)
	err := dec.Decode(&results)
	if err != nil {
		log.Fatalf("failed to unmarshal results: %v", err)
	}

	for _, result := range results {
		showStatsForTest(&result)
	}
}

func showStatsForTest(result *harness.TestResult) {
	fmt.Printf("---\nTest type: %s\nNumber of trials: %d\n", result.TestType, len(result.TrialResults))

	var failures int
	successTimes := make([]time.Duration, 0, len(result.TrialResults))
	allTimes := make([]time.Duration, len(result.TrialResults))
	for i, trial := range result.TrialResults {
		allTimes[i] = trial.Duration
		if trial.Error != "" {
			failures++
		} else {
			successTimes = append(successTimes, trial.Duration)
		}
	}

	successData := stats.LoadRawData(successTimes)
	successMedian, _ := successData.Median()
	successMean, _ := successData.Mean()
	allData := stats.LoadRawData(allTimes)
	allMedian, _ := allData.Median()
	allMean, _ := allData.Mean()

	fmt.Printf("Number of failures: %d\n", failures)
	fmt.Printf("Mean duration for successful trials: %v\n", time.Duration(successMean))
	fmt.Printf("Median duration for successful trials: %v\n\n", time.Duration(successMedian))
	fmt.Printf("Mean duration for all trials: %v\n", time.Duration(allMean))
	fmt.Printf("Median duration for all trials: %v\n\n", time.Duration(allMedian))
}
