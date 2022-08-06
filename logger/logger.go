package logger

import (
	"log"
	"os"
	"tmux_compose/types"
)

func GetStdHandles() types.StdHandlesType {
	return &types.StdHandlesStruct{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

func Construct(stdHandles types.StdHandlesType) Logger {
	logger := Logger{}
	logger.New(stdHandles)
	return logger
}

type Logger struct {
	logger *log.Logger
}

func (logger Logger) New(stdHandles types.StdHandlesType) {
	logger.logger = log.New(stdHandles.Stderr, "", 0)
}

func (logger Logger) Log(s string) {
	logger.logger.Output(2, s)
}
