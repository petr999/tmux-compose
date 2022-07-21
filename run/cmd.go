package run

import (
	"io"
	"os/exec"
)

type execInterface interface {
	GetCommand(name string, arg ...string) Cmd
}

type CmdInterface struct {
	Run func() error
}

type ExecStruct struct{}

func (execStruct ExecStruct) GetCommand(name string, arg ...string) Cmd {
	cmd := exec.Command(name, arg...)
	return Cmd{
		cmd,
		&cmd.Stdout,
		&cmd.Stderr,
		&cmd.Stdin,
	}
}

type Cmd struct {
	Obj interface {
		Run() error
	}
	Stdout *io.Writer
	Stderr *io.Writer
	Stdin  *io.Reader
}

func (cmd Cmd) Run() error {
	return cmd.Obj.Run()
}

func (cmd Cmd) StdCommute(os *OsStruct) error {
	*cmd.Stdout = os.Stdout
	*cmd.Stderr = os.Stderr
	*cmd.Stdin = os.Stdin

	return nil
}

type OsStruct struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}
