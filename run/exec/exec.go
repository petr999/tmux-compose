package exec

import (
	"encoding/json"
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

type cmdObjDryRun struct {
	nameArgs NameArgsType
	stdout   io.Writer
}

func (cmd *cmdObjDryRun) Run() error {
	outVal, err := json.Marshal([]any{cmd.nameArgs.Name, cmd.nameArgs.Args})
	if err != nil {
		log.Panic(`can not serialize cmd name and args`)
	}
	_, err = cmd.stdout.Write(append(outVal, "\n"...))
	if err != nil {
		log.Panicf(`Writing output value: '%v'`, err)
	}
	return nil
}

func (execStruct *ExecStruct) MakeCommand(dryRun *MakeCommandDryRunType,
	nameArgs NameArgsType) {

	var execStructCmd *CmdType
	if len(dryRun.DryRun) > 0 {
		cmd := &cmdObjDryRun{nameArgs, dryRun.OsStruct.Stdout}
		execStructCmd = &CmdType{
			Obj:    cmd,
			Stdout: &dryRun.OsStruct.Stdout,
			Stderr: &dryRun.OsStruct.Stderr,
			Stdin:  &dryRun.OsStruct.Stdin,
		}
	} else {
		cmd := exec.Command(nameArgs.Name, nameArgs.Args...)
		execStructCmd = &CmdType{
			Obj:    cmd,
			Stdout: &cmd.Stdout,
			Stderr: &cmd.Stderr,
			Stdin:  &cmd.Stdin,
		}
	}

	// }
	execStruct.SetCommand(execStructCmd)
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
