package ripcord

import (
	logrustest "github.com/Sirupsen/logrus/hooks/test"
	"testing"
	"time"
)

var (
	logger, hook = logrustest.NewNullLogger()
)

type staticStatsCollector struct {
	static IOStats
}

func (collector *staticStatsCollector) CollectStats() (IOStats, error) {
	return collector.static, nil
}

func TestStaticStats(t *testing.T) {

}
