package ripcord

import "github.com/shirou/gopsutil/net"

// IOStats is a slice of net.IOCounterStat
type IOStats []net.IOCountersStat

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
func (stats IOStats) Find(name string) (stat net.IOCountersStat) {
	for _, s := range stats {
		if s.Name == name {
			return s
		}
	}

	return
}
