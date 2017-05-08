package ripcord

import "fmt"

// ErrBytesRecv is raised when the number of bytes received
// on an interface exceeds the number of allowed bytes.
type ErrBytesRecv struct {
	interfaceName string
	amount        uint64
}

// NewErrBytesRecv returns a pointer to a new instance of a
// ErrBytesRecv error.
func NewErrBytesRecv(interfaceName string, amount uint64) (err *ErrBytesRecv) {
	return &ErrBytesRecv{
		interfaceName: interfaceName,
		amount:        amount,
	}
}

func (err *ErrBytesRecv) Error() string {
	return fmt.Sprintf("%s bytes recv exceeded by %d", err.interfaceName, err.amount)
}

// ErrBytesSent is raised when the number of bytes sent
// on an interface exceeds the number of allowed bytes.
type ErrBytesSent struct {
	interfaceName string
	amount        uint64
}

// NewErrBytesSent returns a pointer to a new instance of a
// ErrBytesSent error.
func NewErrBytesSent(interfaceName string, amount uint64) (err *ErrBytesSent) {
	return &ErrBytesSent{
		interfaceName: interfaceName,
		amount:        amount,
	}
}

func (err *ErrBytesSent) Error() string {
	return fmt.Sprintf("%s bytes sent exceeded by %d", err.interfaceName, err.amount)
}
