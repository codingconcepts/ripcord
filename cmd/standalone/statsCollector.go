package main

import (
	"github.com/codingconcepts/ripcord"
	"github.com/shirou/gopsutil/net"
)

// netStatsCollector is a StatsCollector implementation
// which retrieve real network statistics.
type netStatsCollector struct {
}

// CollectStats collects real network statistics.
func (c *netStatsCollector) CollectStats() (stats ripcord.IOStats, err error) {
	var netStats []net.IOCountersStat
	if netStats, err = net.IOCounters(true); err != nil {
		return nil, err
	}

	stats = make(ripcord.IOStats, len(netStats))
	for i, ns := range netStats {
		stats[i] = ripcord.IOStat{
			Name:      ns.Name,
			BytesRecv: ns.BytesRecv,
			BytesSent: ns.BytesSent,
		}
	}

	return
}
