package harness

import (
	"log"
	"sync"
	"time"

	"github.com/adamwg/distribution-bench/config"
	"github.com/adamwg/distribution-bench/tests"
	"github.com/bloodorangeio/reggie"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Run runs a benchmark based on the given config.
func Run(cfg *config.Config) ([]TestResult, error) {
	client, err := reggie.NewClient(
		cfg.Registry.RootURL,
		reggie.WithUserAgent("distribution-bench"),
		reggie.WithUsernamePassword(cfg.Registry.Username, cfg.Registry.Password),
		reggie.WithDefaultName(cfg.Registry.Namespace),
		reggie.WithDebug(cfg.Debug),
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating registry client")
	}

	results := make([]TestResult, len(cfg.Tests))
	for i, test := range cfg.Tests {
		result, err := runTest(client, test)
		if err != nil {
			return nil, errors.Wrap(err, "running test")
		}
		results[i] = result
	}

	return results, nil
}

func runTest(client *reggie.Client, cfg config.TestConfig) (TestResult, error) {
	log.Printf("Starting %s test run with %d trials and concurrency %d",
		cfg.Type, cfg.Trials, cfg.Concurrency)

	result := TestResult{
		TestType:     cfg.Type,
		TrialResults: make([]TrialResult, cfg.Trials),
	}

	trials := make([]tests.Runner, cfg.Trials)
	for i := range trials {
		trial, err := tests.CreateRunner(cfg)
		if err != nil {
			return result, errors.Wrap(err, "creating runner for test")
		}
		trials[i] = trial
	}

	var eg errgroup.Group
	for _, t := range trials {
		t := t
		eg.Go(func() error {
			return t.Prepare(client)
		})
	}
	err := eg.Wait()
	if err != nil {
		return result, errors.Wrap(err, "preparing trials")
	}

	runCh := make(chan tests.Runner, cfg.Trials)
	resultCh := make(chan TrialResult, cfg.Trials)

	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go testWorker(client, runCh, resultCh, &wg)
	}

	for _, t := range trials {
		runCh <- t
	}
	close(runCh)

	for i := range trials {
		result.TrialResults[i] = <-resultCh
	}

	wg.Wait()

	return result, nil
}

func testWorker(client *reggie.Client, runCh <-chan tests.Runner, resultCh chan<- TrialResult, wg *sync.WaitGroup) {
	var (
		start    time.Time
		duration time.Duration
		bytes    uint64
		err      error
	)

	for trial := range runCh {
		log.Print("starting trial")

		start = time.Now()
		bytes, err = trial.Run(client)
		duration = time.Since(start)

		res := TrialResult{
			Bytes:    bytes,
			Duration: duration,
		}
		if err != nil {
			res.Error = err.Error()
		}

		log.Print("trial completed")
		resultCh <- res
	}

	wg.Done()
}
