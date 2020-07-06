package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/adamwg/distribution-bench/config"
	"github.com/adamwg/distribution-bench/harness"
)

func main() {
	configFilePath := "-"

	if len(os.Args) >= 2 {
		configFilePath = os.Args[1]
	}

	var configFile *os.File
	if configFilePath == "-" {
		configFile = os.Stdin
	} else {
		fd, err := os.Open(configFilePath)
		if err != nil {
			log.Fatalf("failed to open config file %q: %v", configFilePath, err)
		}
		defer fd.Close()
		configFile = fd
	}

	var cfg config.Config
	dec := json.NewDecoder(configFile)
	err := dec.Decode(&cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	results, err := harness.Run(&cfg)
	if err != nil {
		log.Fatalf("running tests failed: %v", err)
	}
	err = json.NewEncoder(os.Stdout).Encode(results)
	if err != nil {
		log.Fatalf("encoding results: %v", err)
	}
}
