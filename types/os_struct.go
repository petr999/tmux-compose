package types

import (
	"io"
	"os/exec"
)

type StdHandlesStruct struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

type StdHandlesType *StdHandlesStruct

type ExecOsInterface interface {
	Getenv(key string) string
	Chdir(dir string) error
	ReadFile(name string) ([]byte, error)
	GetStdHandles() StdHandlesType
	Command(name string, arg ...string) *exec.Cmd
}

type RunnerOsInterface interface {
	Exit(code int)
}
