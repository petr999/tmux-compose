package exec

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

type ExecInterface interface {
	GetCommand(name string, arg ...string) cmdType
}

type CmdInterface struct {
	Run func() error
}

type ExecStruct struct{}

func (execStruct ExecStruct) GetCommand(name string, arg ...string) cmdType {
	cmd := exec.Command(name, arg...)
	return cmdType{
		cmd,
		&cmd.Stdout,
		&cmd.Stderr,
		&cmd.Stdin,
	}
}

type cmdType struct {
	Obj interface {
		Run() error
	}
	Stdout *io.Writer
	Stderr *io.Writer
	Stdin  *io.Reader
}

func (cmd cmdType) Run() error {
	return cmd.Obj.Run()
}

func (cmd cmdType) StdCommute(os *OsStruct) error {
	*cmd.Stdout = os.Stdout
	*cmd.Stderr = os.Stderr
	*cmd.Stdin = os.Stdin

	return nil
}

type OsStruct struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	Exit   func(code int)
}

func LogFunc(v ...any) {
	log.Output(2, fmt.Sprint(v...))
}
