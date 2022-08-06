package types

import (
	"io"
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
	Getwd() (dir string, err error)
	ReadFile(name string) ([]byte, error)
	GetStdHandles() StdHandlesType
}

type RunnerOsInterface interface {
	Exit(code int)
}
