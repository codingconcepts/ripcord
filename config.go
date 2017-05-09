package ripcord

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"
)

// InterfaceConfigs defines a set of InterfaceConfig settings,
// allowing you to react to different thresholds per network
// interface.
type InterfaceConfigs struct {
	CheckInterval ConfigDuration    `json:"checkInterval"`
	Interfaces    []InterfaceConfig `json:"interfaces"`
}

// InterfaceConfig holds the configurable thresholds for a
// network interface.
type InterfaceConfig struct {
	Name         string   `json:"name"`
	MaxBytesRecv uint64   `json:"maxBytesRecv"`
	MaxBytesSent uint64   `json:"maxBytesSent"`
	Instructions []string `json:"instructions"`
}

// ConfigDuration allows for the configuration of durations
// in the form of "5m30s", as opposed to the default Unix
// epoch timestamp.
type ConfigDuration struct {
	time.Duration
}

// NewConfigsFromReader returns a pointer to a new instance of an
// InterfaceConfigs struct from a reader.
func NewConfigsFromReader(reader io.Reader) (configs InterfaceConfigs, err error) {
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, reader); err != nil {
		return
	}

	configs = InterfaceConfigs{}
	err = json.Unmarshal(buf.Bytes(), &configs)

	return
}

// CompareStats compares the last snapshot of IOStats
// against the current snapshot of IOStats and returns
// an error if any of the thresholds have been breached.
func (configs InterfaceConfigs) CompareStats(prev IOStats, curr IOStats) (err error) {
	for _, config := range configs.Interfaces {
		p := prev.Find(config.Name)
		c := curr.Find(config.Name)

		if err = config.CompareStat(p, c); err != nil {
			return
		}
	}

	return
}

// CompareStat compares the last snapshot of an individual
// IOStat against the current snapshot of the same IOStats
// and returns an error if any of the thresholds have been
// breached.
func (config InterfaceConfig) CompareStat(prev IOStat, curr IOStat) (err error) {
	bytesRecvDiff := curr.BytesRecv - prev.BytesRecv
	if curr.BytesRecv > prev.BytesRecv && bytesRecvDiff > config.MaxBytesRecv {
		return NewErrBytesRecv(prev.Name, bytesRecvDiff, config)
	}

	bytesSentDiff := curr.BytesSent - prev.BytesSent
	if curr.BytesSent > prev.BytesSent && bytesSentDiff > config.MaxBytesSent {
		return NewErrBytesSent(prev.Name, bytesSentDiff, config)
	}

	return
}

// UnmarshalJSON unmarshals a ConfigDuration from JSON.
func (d *ConfigDuration) UnmarshalJSON(b []byte) (err error) {
	// trim off the quotes to get at the duration (the
	// JSON object will appear as "1s")
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}
