package exec

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

type ExecInterface interface {
	MakeCommand(name string, arg ...string)
	GetCommand() *cmdType
}

type CmdInterface struct {
	Run func() error
}

type ExecStruct struct {
	cmd *cmdType
}

func (execStruct *ExecStruct) MakeCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	execStruct.cmd = &cmdType{
		cmd,
		&cmd.Stdout,
		&cmd.Stderr,
		&cmd.Stdin,
	}
}

func (execStruct *ExecStruct) GetCommand() *cmdType {
	return execStruct.cmd
}

type cmdType struct {
	Obj interface {
		Run() error
	}
	Stdout *io.Writer
	Stderr *io.Writer
	Stdin  *io.Reader
}

func (cmd *cmdType) Run() error {
	return cmd.Obj.Run()
}

func (cmd *cmdType) StdCommute(os *OsStruct) error {
	*cmd.Stdout = os.Stdout
	*cmd.Stderr = os.Stderr
	*cmd.Stdin = os.Stdin

	return nil
}

type OsStructExit func(code int)

type OsStruct struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	Exit   OsStructExit
}

func LogFunc(v ...any) {
	log.Output(2, fmt.Sprint(v...))
}
