package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codingconcepts/ripcord"
)

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		panic("handle")
	}

	configs, err := ripcord.NewConfigsFromReader(file)
	if err != nil {
		panic("handle")
	}

	logger := ripcord.NewLogger(os.Stdout, log.DebugLevel)
	collector := new(netStatsCollector)
	runner := ripcord.NewRunner(collector, configs, logger)

	go func() {
		if err := runner.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	fmt.Scanln()
	runner.Stop()
}
