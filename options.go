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
