package ripcord

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Runner contains all of the properties necessary for
// running an instance of Ripcord, standalone or otherwise.
type Runner struct {
	StatsCollector StatsCollector
	configs        InterfaceConfigs
	logger         *log.Logger
}

// NewRunner returns the pointer to a new instance of
// a Ripcord struct.
func NewRunner(statsCollector StatsCollector, configs InterfaceConfigs, logger *log.Logger) (runner *Runner) {
	return &Runner{
		StatsCollector: statsCollector,
		configs:        configs,
		logger:         logger,
	}
}

// Start kicks everything off.  It blocks, so launch in
// separate goroutine.
func (runner *Runner) Start(configs InterfaceConfigs) (err error) {
	var prev IOStats
	var curr IOStats

	if prev, err = runner.StatsCollector.CollectStats(); err != nil {
		return
	}

	for {
		select {
		case <-time.Tick(runner.configs.CheckInterval.Duration):
			if curr, err = runner.StatsCollector.CollectStats(); err != nil {
				return
			}

			if err = runner.configs.CompareStats(prev, curr); err != nil {
				runner.logger.WithError(err).Error("breach detected")

				// if the error encounterd satisfies the command executor
				// interface, execute its instructions now
				if executor, ok := err.(CommandExecutor); ok {
					runner.logger.WithError(err).Error("executing command")
					return executor.Execute()
				}
			} else {
				runner.logger.Debug("no breach detected")
			}

			prev = curr
		}
	}
}
