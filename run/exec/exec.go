package exec

import (
	"io"
	"log"
	"os/exec"
)

type OsStructInterface interface {
	Exit(code int)
	Getenv(key string) string
}

type MakeCommandDryRunType struct {
	DryRun   string
	OsStruct *OsStruct
}
type NameArgsType struct {
	Name string
	Args []string
}

type ExecInterface interface {
	MakeCommand(*MakeCommandDryRunType,
		NameArgsType,
	)
	GetCommand() *CmdType
	SetCommand(*CmdType)
}

type CmdInterface struct {
	Run func() error
}

type ExecStruct struct {
	cmd *CmdType
}

type dryRunCommand struct {
}

func (cmd *dryRunCommand) Run() error {
	return nil
}

func (execStruct *ExecStruct) MakeCommand(dryRun *MakeCommandDryRunType,
	nameArgs NameArgsType) {
	// cmd & CmdType{
	// 	&cmd.Stdout,
	// 	&cmd.Stderr,
	// 	&cmd.Stdin,
	// }

	// func(name, arg...){

	// }
	// if len(dryRun) == 0 {
	cmd := exec.Command(nameArgs.Name, nameArgs.Args...)
	// }
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
