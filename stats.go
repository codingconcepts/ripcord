package ripcord

// StatsCollector defines the behaviour of something
// which can collect stats.
type StatsCollector interface {
	CollectStats() (IOStats, error)
}

// IOStat is represents a snapshot for a network interface.
type IOStat struct {
	Name      string
	BytesSent uint64
	BytesRecv uint64
}

// IOStats is a slice of IOStat structs.
type IOStats []IOStat

// Filter filters a given IOStats agains a set of interface
// names, returning only those whose names match.
func (stats IOStats) Filter(names ...string) (filtered IOStats) {
	filtered = stats[:0]
	for _, name := range names {
		stat := stats.Find(name)
		if stat.Name != "" {
			filtered = append(filtered, stat)
		}
	}
	return
}

// Find returns the first net.IOCounterStat whose name
// matches the given value.
func (stats IOStats) Find(name string) (stat IOStat) {
	for _, s := range stats {
		if s.Name == name {
			return s
		}
	}

	return
}
