package ripcord

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

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
		`{
			"checkInterval": "5s",
			"interfaces": [
				{
					"name": "WiFi",
					"maxBytesRecv": 123,
					"maxBytesSent": 234,
					"instructions": [ "echo", "WiFi breached" ]
				},
				{
					"name": "Local Area Connection *1",
					"maxBytesRecv": 345,
					"maxBytesSent": 456,
					"instructions": [ "echo", "Local Area Connection *1 breached" ]
				}
			]
		}`)

	configs, err := NewConfigsFromReader(reader)
	test.ErrorNil(t, err)

	test.Equals(t, time.Second*5, configs.CheckInterval.Duration)
	test.Equals(t, 2, len(configs.Interfaces))

	test.Equals(t, "WiFi", configs.Interfaces[0].Name)
	test.Equals(t, uint64(123), configs.Interfaces[0].MaxBytesRecv)
	test.Equals(t, uint64(234), configs.Interfaces[0].MaxBytesSent)
	test.Equals(t, []string{"echo", "WiFi breached"}, configs.Interfaces[0].Instructions)

	test.Equals(t, "Local Area Connection *1", configs.Interfaces[1].Name)
	test.Equals(t, uint64(345), configs.Interfaces[1].MaxBytesRecv)
	test.Equals(t, uint64(456), configs.Interfaces[1].MaxBytesSent)
	test.Equals(t, []string{"echo", "Local Area Connection *1 breached"}, configs.Interfaces[1].Instructions)
}

func TestNewConfigFromReaderReturnsError(t *testing.T) {
	configs, err := NewConfigsFromReader(new(errorReader))
	test.ErrorNotNil(t, err)
	test.Equals(t, errFromErrorReader, err)
	test.Equals(t, 0, len(configs.Interfaces))
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

func TestCompareStatReturnsErrorIfStatsExceedThreshold(t *testing.T) {
	config := &InterfaceConfig{Name: "a", MaxBytesRecv: 10}
	err := config.CompareStat(IOStat{Name: "a", BytesRecv: 10}, IOStat{Name: "a", BytesRecv: 21})
	test.ErrorNotNil(t, err)

	recvErr, ok := err.(*ErrBytesRecv)
	test.Assert(t, ok)
	test.Equals(t, "a", recvErr.interfaceName)
	test.Equals(t, uint64(11), recvErr.amount)

	config = &InterfaceConfig{Name: "b", MaxBytesSent: 10}
	err = config.CompareStat(IOStat{Name: "b", BytesSent: 10}, IOStat{Name: "b", BytesSent: 21})
	test.ErrorNotNil(t, err)

	sentErr, ok := err.(*ErrBytesSent)
	test.Assert(t, ok)
	test.Equals(t, "b", sentErr.interfaceName)
	test.Equals(t, uint64(11), sentErr.amount)
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

func TestCompareStatsReturnNoErrorIfBytesRecvWithinThreshold(t *testing.T) {
	configs := InterfaceConfigs{
		Interfaces: []InterfaceConfig{
			InterfaceConfig{
				Name: "a", MaxBytesRecv: 10, MaxBytesSent: 10},
		},
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
		Interfaces: []InterfaceConfig{
			InterfaceConfig{
				Name: "a", MaxBytesRecv: 10, MaxBytesSent: 10},
		},
	}

	err := configs.CompareStats(IOStats{IOStat{Name: "a", BytesRecv: 10}}, IOStats{IOStat{Name: "a", BytesRecv: 21}})
	test.ErrorNotNil(t, err)

	recvErr, ok := err.(*ErrBytesRecv)
	test.Assert(t, ok)
	test.Equals(t, "a", recvErr.interfaceName)
	test.Equals(t, uint64(11), recvErr.amount)

	err = configs.CompareStats(IOStats{IOStat{Name: "a", BytesSent: 10}}, IOStats{IOStat{Name: "a", BytesSent: 21}})
	test.ErrorNotNil(t, err)

	sentErr, ok := err.(*ErrBytesSent)
	test.Assert(t, ok)
	test.Equals(t, "a", sentErr.interfaceName)
	test.Equals(t, uint64(11), sentErr.amount)
}

func TestConfigDuration(t *testing.T) {
	var s struct {
		Duration ConfigDuration `json:"duration"`
	}

	j := `{ "duration": "1h30m40s" }`

	err := json.Unmarshal([]byte(j), &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.Duration.Duration != time.Hour+time.Minute*30+time.Second*40 {
		t.Fatalf("expected 1h30m40s but got %v", s.Duration)
	}
}
