package harness

import (
	"time"

	"github.com/adamwg/distribution-bench/config"
)

// TrialResult is the result of a single test trial.
type TrialResult struct {
	Bytes    uint64        `json:"bytes"`
	Duration time.Duration `json:"duration_nanos"`
	Error    string        `json:"error"`
}

// TestResult is the result of a set of trials.
type TestResult struct {
	TestType     config.TestType `json:"test_type"`
	TrialResults []TrialResult   `json:"trial_results"`
}
