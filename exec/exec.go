package exec

import (
	osExec "os/exec"
	"tmux_compose/types"
)

func Construct(osStruct types.ExecOsInterface) *Exec {
	exec := &Exec{}
	stdHandles := osStruct.GetStdHandles()
	exec.New(osStruct, stdHandles)
	return exec
}

type dryRunCmd struct{}

func (obj dryRunCmd) Run() error { return nil }

type Exec struct {
	Cmd        *types.CmdType
	Selector   any
	osStruct   *types.ExecOsInterface
	stdHandles types.StdHandlesType
}

func (exec *Exec) New(osStruct types.ExecOsInterface, stdHandles types.StdHandlesType) {
	exec.Selector, exec.osStruct, exec.stdHandles = false, &osStruct, stdHandles
}

func (exec *Exec) execCommand(cna types.CmdNameArgsValueType) *osExec.Cmd {
	obj := osExec.Command(cna.Name, cna.Args...)
	obj.Stdout = exec.stdHandles.Stdout
	obj.Stderr = exec.stdHandles.Stderr
	obj.Stdin = exec.stdHandles.Stdin

	return obj
}

func (exec *Exec) dryRun(cna types.CmdNameArgsValueType) types.CmdInterface {
	return dryRunCmd{}
}

func (exec *Exec) getSelector() any {
	return false
}

func (exec *Exec) GetCommand(cna types.CmdNameArgsValueType) *types.CmdType {

	var obj interface{ Run() error }
	selector := exec.getSelector()
	if dryRun, ok := selector.(bool); ok {
		if dryRun {
			obj = exec.dryRun(cna)
		} else {
			obj = exec.execCommand(cna)
		}
	} else {
		if objWithRun, ok := selector.(interface{ Run() error }); ok {
			obj = objWithRun
		}
	}
	exec.Cmd = &types.CmdType{Obj: obj, Stdhandles: exec.stdHandles}
	return exec.Cmd
}

// type ExecInterface interface {
// 	MakeCommand(*MakeCommandDryRunType,
// 		nameArgsType,
// 	)
// 	GetCommand() *CmdType
// 	SetCommand(*CmdType)
// }

// type CmdInterface struct {
// 	Run func() error
// }

// type OsStructInterface interface {
// 	Exit(code int)
// 	Getenv(key string) string
// 	Chdir(dir string) error
// }

// type MakeCommandDryRunType struct {
// 	DryRun   string
// 	OsStruct *types.OsStruct
// }
// type nameArgsType = types.CmdNameArgsType

// type ExecStruct struct {
// 	cmd *CmdType
// }

// type cmdObjDryRun struct {
// 	nameArgs nameArgsType
// 	stdout   io.Writer
// }

// func (cmd *cmdObjDryRun) Run() error {
// 	outVal, err := json.Marshal([]any{cmd.nameArgs.Name, cmd.nameArgs.Args})
// 	if err != nil {
// 		log.Panic(`can not serialize cmd name and args`)
// 	}
// 	_, err = cmd.stdout.Write(append(outVal, "\n"...))
// 	if err != nil {
// 		log.Panicf(`Writing output value: '%v'`, err)
// 	}
// 	return nil
// }

// func (execStruct *ExecStruct) MakeCommand(dryRun *MakeCommandDryRunType,
// 	nameArgs nameArgsType) {

// 	var execStructCmd *CmdType
// 	if len(dryRun.DryRun) > 0 {
// 		cmd := &cmdObjDryRun{nameArgs, dryRun.OsStruct.Stdout}
// 		execStructCmd = &CmdType{
// 			Obj:    cmd,
// 			Stdout: &dryRun.OsStruct.Stdout,
// 			Stderr: &dryRun.OsStruct.Stderr,
// 			Stdin:  &dryRun.OsStruct.Stdin,
// 		}
// 	} else {
// 		cmd := exec.Command(nameArgs.Name, nameArgs.Args...)
// 		execStructCmd = &CmdType{
// 			Obj:    cmd,
// 			Stdout: &cmd.Stdout,
// 			Stderr: &cmd.Stderr,
// 			Stdin:  &cmd.Stdin,
// 		}
// 	}

// 	// }
// 	execStruct.SetCommand(execStructCmd)
// }

// func (execStruct *ExecStruct) GetCommand() *CmdType {
// 	return execStruct.cmd
// }

// func (execStruct *ExecStruct) SetCommand(cmd *CmdType) {
// 	execStruct.cmd = cmd
// }

// type CmdType struct {
// 	Obj interface {
// 		Run() error
// 	}
// 	Stdout *io.Writer
// 	Stderr *io.Writer
// 	Stdin  *io.Reader
// }

// func (cmd *CmdType) Run() error {
// 	return cmd.Obj.Run()
// }

// func (Cmd *CmdType) StdCommute(os *types.OsStruct) error {
// 	*Cmd.Stdout = os.Stdout
// 	*Cmd.Stderr = os.Stderr
// 	*Cmd.Stdin = os.Stdin

// 	return nil
// }
