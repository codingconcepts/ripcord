package ripcord

import (
	"errors"
	"strings"
	"testing"

	"github.com/codingconcepts/ripcord/test"
)

var (
	errFromErrorReader = errors.New("error loading config")
)

type errorReader struct{}

func (r *errorReader) Read(data []byte) (n int, err error) {
	return 0, errFromErrorReader
}

func TestNewConfigFromReader(t *testing.T) {
	reader := strings.NewReader(
		`[
			{
				"name": "WiFi",
				"maxBytesRecv": 500000,
				"maxBytesSent": 1000000
			},
			{
				"name": "Local Area Connection *1",
				"maxBytesRecv": 20000000,
				"maxBytesSent": 30000000
			}
		]`)

	configs, err := NewConfigsFromReader(reader)
	test.ErrorNil(t, err)
	test.Equals(t, 2, len(configs))

	test.Equals(t, "WiFi", configs[0].Name)
	test.Equals(t, uint64(500000), configs[0].MaxBytesRecv)
	test.Equals(t, uint64(1000000), configs[0].MaxBytesSent)

	test.Equals(t, "Local Area Connection *1", configs[1].Name)
	test.Equals(t, uint64(20000000), configs[1].MaxBytesRecv)
	test.Equals(t, uint64(30000000), configs[1].MaxBytesSent)
}

func TestNewConfigFromReaderReturnsError(t *testing.T) {
	configs, err := NewConfigsFromReader(new(errorReader))
	test.ErrorNotNil(t, err)
	test.Equals(t, errFromErrorReader, err)
	test.Equals(t, 0, len(configs))
}

func TestCompareStatReturnNoErrorIfBytesRecvWithinThreshold(t *testing.T) {
	config := &InterfaceConfig{
		Name:         "a",
		MaxBytesRecv: 10,
	}

	test.ErrorNil(t, config.CompareStat(
		IOStat{Name: "a", BytesRecv: 10},
		IOStat{Name: "a", BytesRecv: 20}))
}

func TestCompareStatReturnsErrorIfBytesRecvExceedsThreshold(t *testing.T) {
	config := &InterfaceConfig{
		Name:         "a",
		MaxBytesRecv: 10,
	}

	err := config.CompareStat(
		IOStat{Name: "a", BytesRecv: 10},
		IOStat{Name: "a", BytesRecv: 21})

	test.ErrorNotNil(t, err)

	realErr, ok := err.(*ErrBytesRecv)
	test.Assert(t, ok)
	test.Equals(t, "a", realErr.interfaceName)
	test.Equals(t, uint64(11), realErr.amount)
}

func TestCompareStatReturnNoErrorIfBytesSentWithinThreshold(t *testing.T) {
	config := &InterfaceConfig{
		Name:         "a",
		MaxBytesSent: 10,
	}

	test.ErrorNil(t, config.CompareStat(
		IOStat{Name: "a", BytesSent: 10},
		IOStat{Name: "a", BytesSent: 20}))
}

func TestCompareStatReturnsErrorIfBytesSentExceedsThreshold(t *testing.T) {
	config := &InterfaceConfig{
		Name:         "a",
		MaxBytesSent: 10,
	}

	err := config.CompareStat(
		IOStat{Name: "a", BytesSent: 10},
		IOStat{Name: "a", BytesSent: 21})

	test.ErrorNotNil(t, err)

	realErr, ok := err.(*ErrBytesSent)
	test.Assert(t, ok)
	test.Equals(t, "a", realErr.interfaceName)
	test.Equals(t, uint64(11), realErr.amount)
}

func TestCompareStatsReturnNoErrorIfBytesRecvWithinThreshold(t *testing.T) {
	configs := InterfaceConfigs{
		InterfaceConfig{
			Name: "a", MaxBytesRecv: 10, MaxBytesSent: 10},
	}

	prev := IOStats{
		IOStat{Name: "a", BytesRecv: 10, BytesSent: 10},
	}

	curr := IOStats{
		IOStat{Name: "a", BytesRecv: 20, BytesSent: 20},
	}

	test.ErrorNil(t, configs.CompareStats(prev, curr))
}

func TestCompareStatsReturnsErrorIfBytesRecvExceedsThreshold(t *testing.T) {
	configs := InterfaceConfigs{
		InterfaceConfig{Name: "a", MaxBytesRecv: 10},
	}

	prev := IOStats{IOStat{Name: "a", BytesRecv: 10}}
	curr := IOStats{IOStat{Name: "a", BytesRecv: 21}}

	err := configs.CompareStats(prev, curr)
	test.ErrorNotNil(t, err)

	realErr, ok := err.(*ErrBytesRecv)
	test.Assert(t, ok)
	test.Equals(t, "a", realErr.interfaceName)
	test.Equals(t, uint64(11), realErr.amount)
}

func TestCompareStatsReturnsErrorIfBytesSentExceedsThreshold(t *testing.T) {
	configs := InterfaceConfigs{
		InterfaceConfig{Name: "a", MaxBytesRecv: 10},
	}

	prev := IOStats{IOStat{Name: "a", BytesSent: 10}}
	curr := IOStats{IOStat{Name: "a", BytesSent: 21}}

	err := configs.CompareStats(prev, curr)
	test.ErrorNotNil(t, err)

	realErr, ok := err.(*ErrBytesSent)
	test.Assert(t, ok)
	test.Equals(t, "a", realErr.interfaceName)
	test.Equals(t, uint64(11), realErr.amount)
}
