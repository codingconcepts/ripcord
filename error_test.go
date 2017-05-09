package ripcord

import (
	"testing"

	"github.com/codingconcepts/ripcord/test"
)

func TestErrBytesRecvProvidesCorrectMessage(t *testing.T) {
	err := ErrBytesRecv{
		amount:        1234,
		interfaceName: "WiFi",
	}

	test.Equals(t, "WiFi bytes recv exceeded by 1234", err.Error())
}

func TestErrBytesSentProvidesCorrectMessage(t *testing.T) {
	err := ErrBytesSent{
		amount:        4321,
		interfaceName: "Thing",
	}

	test.Equals(t, "Thing bytes sent exceeded by 4321", err.Error())
}
