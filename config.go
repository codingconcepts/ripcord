package ripcord

import (
	"bytes"
	"encoding/json"
	"io"
)

// InterfaceConfigs defines a set of InterfaceConfig settings,
// allowing you to react to different thresholds per network
// interface.
type InterfaceConfigs []InterfaceConfig

// InterfaceConfig holds the configurable thresholds for a
// network interface.
type InterfaceConfig struct {
	Name         string
	MaxBytesRecv uint64
	MaxBytesSent uint64
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
	for _, config := range configs {
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
		return NewErrBytesRecv(prev.Name, bytesRecvDiff)
	}

	bytesSentDiff := curr.BytesSent - prev.BytesSent
	if curr.BytesSent > prev.BytesSent && bytesSentDiff > config.MaxBytesSent {
		return NewErrBytesSent(prev.Name, bytesSentDiff)
	}

	return
}
