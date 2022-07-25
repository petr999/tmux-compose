package exec

import (
	"io"
	"log"
	"os/exec"
)

type ExecInterface interface {
	MakeCommand(name string, arg ...string)
	GetCommand() *CmdType
	SetCommand(*CmdType)
}

type CmdInterface struct {
	Run func() error
}

type ExecStruct struct {
	cmd *CmdType
}

func (execStruct *ExecStruct) MakeCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	execStruct.SetCommand(&CmdType{
		cmd,
		&cmd.Stdout,
		&cmd.Stderr,
		&cmd.Stdin,
	})
}

func (execStruct *ExecStruct) GetCommand() *CmdType {
	return execStruct.cmd
}

func (execStruct *ExecStruct) SetCommand(cmd *CmdType) {
	execStruct.cmd = cmd
}

type CmdType struct {
	Obj interface {
		Run() error
	}
	Stdout *io.Writer
	Stderr *io.Writer
	Stdin  *io.Reader
}

func (cmd *CmdType) Run() error {
	return cmd.Obj.Run()
}

func (Cmd *CmdType) StdCommute(os *OsStruct) error {
	*Cmd.Stdout = os.Stdout
	*Cmd.Stderr = os.Stderr
	*Cmd.Stdin = os.Stdin

	return nil
}

type OsStructExit func(code int)
type OsStructGetEnv func(key string) string

type OsStruct struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	Exit   OsStructExit
	Getenv OsStructGetEnv
}

func GetLogFunc(writer io.Writer) func(s string) {
	logger := log.New(writer, "", 0)
	return func(s string) {
		logger.Output(2, s)
	}
}
