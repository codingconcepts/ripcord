package ripcord

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/net"
)

// Runner contains all of the properties necessary for
// running an instance of Ripcord, standalone or otherwise.
type Runner struct {
	configs InterfaceConfigs
	logger  *log.Logger
}

// NewRunner returns the pointer to a new instance of
// a Ripcord struct.
func NewRunner(configs InterfaceConfigs, logger *log.Logger) (runner *Runner) {
	return &Runner{
		configs: configs,
		logger:  logger,
	}
}

// Start blocks, launch in separate goroutine
func (runner *Runner) Start(configs InterfaceConfigs) (err error) {
	var prev IOStats
	var curr IOStats

	if prev, err = runner.LoadStats(); err != nil {
		return
	}

	for {
		select {
		case <-time.Tick(time.Second):
			if curr, err = runner.LoadStats(); err != nil {
				return
			}

			if err = configs.compareStats(prev, curr); err != nil {
				runner.logger.WithError(err).Error("breach detected")
			} else {
				runner.logger.Debug("no breach detected")
			}

			prev = curr
		}
	}
}

// LoadStats captures net.IOCounters and returns an error
// if they can't be obtained.
func (runner *Runner) LoadStats() (IOStats, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	return IOStats(stats), nil
}

func (configs InterfaceConfigs) compareStats(prev IOStats, curr IOStats) (err error) {
	for _, config := range configs {
		p := prev.Find(config.Name)
		c := curr.Find(config.Name)

		if err = config.compareStat(p, c); err != nil {
			return
		}
	}

	return
}

func (c InterfaceConfig) compareStat(prev net.IOCountersStat, curr net.IOCountersStat) (err error) {
	bytesRecvDiff := curr.BytesRecv - prev.BytesRecv
	if curr.BytesRecv > prev.BytesRecv && bytesRecvDiff > c.MaxBytesRecv {
		return NewErrBytesRecv(prev.Name, bytesRecvDiff)
	}

	bytesSentDiff := curr.BytesSent - prev.BytesSent
	if curr.BytesSent > prev.BytesSent && bytesSentDiff > c.MaxBytesSent {
		return NewErrBytesSent(prev.Name, bytesSentDiff)
	}

	return
}
