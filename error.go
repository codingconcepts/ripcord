package ripcord

import (
	"fmt"
	"os"
	"os/exec"
)

// CommandExecutor describes the behaviour of something capable
// of executing a command.
type CommandExecutor interface {
	Execute() error
}

// ErrBytesRecv is raised when the number of bytes received
// on an interface exceeds the number of allowed bytes.
type ErrBytesRecv struct {
	interfaceName string
	amount        uint64
	config        InterfaceConfig
}

// NewErrBytesRecv returns a pointer to a new instance of a
// ErrBytesRecv error.
func NewErrBytesRecv(interfaceName string, amount uint64, config InterfaceConfig) (err *ErrBytesRecv) {
	return &ErrBytesRecv{
		interfaceName: interfaceName,
		amount:        amount,
		config:        config,
	}
}

func (err *ErrBytesRecv) Error() string {
	return fmt.Sprintf("%s bytes recv exceeded by %d", err.interfaceName, err.amount)
}

// Execute uses information contained within the error to
// execute a command.
func (err *ErrBytesRecv) Execute() error {
	return executeCommand(err.config)
}

// ErrBytesSent is raised when the number of bytes sent
// on an interface exceeds the number of allowed bytes.
type ErrBytesSent struct {
	interfaceName string
	amount        uint64
	config        InterfaceConfig
}

// NewErrBytesSent returns a pointer to a new instance of a
// ErrBytesSent error.
func NewErrBytesSent(interfaceName string, amount uint64, config InterfaceConfig) (err *ErrBytesSent) {
	return &ErrBytesSent{
		interfaceName: interfaceName,
		amount:        amount,
		config:        config,
	}
}

func (err *ErrBytesSent) Error() string {
	return fmt.Sprintf("%s bytes sent exceeded by %d", err.interfaceName, err.amount)
}

// Execute uses information contained within the error to
// execute a command.
func (err *ErrBytesSent) Execute() error {
	return executeCommand(err.config)
}

func executeCommand(config InterfaceConfig) (err error) {
	name := config.Instructions[0]

	// allow for single commands and commands with arguments to
	// be executed (1: will fail for single commands).
	var args []string
	if len(config.Instructions) > 1 {
		args = config.Instructions[1:]
	}

	fmt.Println(name, args)

	// #nosec - the purpose of this killer is to execute
	// subprocesses with variables...
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
