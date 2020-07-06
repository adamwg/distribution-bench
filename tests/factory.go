package tests

import (
	"errors"
	"sync"

	"github.com/adamwg/distribution-bench/config"
)

type creatorFn func(cfg config.TestConfig) (Runner, error)

var (
	// ErrInvalidType is returned by CreateTest if the test type in the config
	// has no creator.
	ErrInvalidType = errors.New("invalid test type")
	creators       map[config.TestType]creatorFn
	creatorsMu     sync.Mutex
)

// CreateRunner creates a test from the given config.
func CreateRunner(cfg config.TestConfig) (Runner, error) {
	creator, ok := creators[cfg.Type]
	if !ok {
		return nil, ErrInvalidType
	}

	return creator(cfg)
}

func registerCreator(t config.TestType, creator creatorFn) {
	creatorsMu.Lock()
	if creators == nil {
		creators = make(map[config.TestType]creatorFn)
	}
	creators[t] = creator
	creatorsMu.Unlock()
}
