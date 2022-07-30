package run

import (
	"fmt"
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
)

const DryRunEnvVarName = `TMUX_COMPOSE_DRY_RUN`

type LogFuncType func(s string)
type CmdNameArgsType func(dcConfigReader dc_config.ReaderInterface) (string, []string)

type Runner struct {
	CmdNameArgs    CmdNameArgsType
	DcConfigReader dc_config.ReaderInterface
	ExecStruct     exec.ExecInterface
	OsStruct       *exec.OsStruct
	LogFunc        LogFuncType
}

func (runner *Runner) Run() {
	DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	CmdNameArgs := runner.CmdNameArgs
	cmdName, args := CmdNameArgs(DcConfigReader)

	ExecStruct.MakeCommand(&exec.MakeCommandDryRunType{DryRun: OsStruct.Getenv(DryRunEnvVarName), OsStruct: OsStruct},
		exec.NameArgsType{Name: cmdName, Args: args})
	cmd := ExecStruct.GetCommand()
	cmd.StdCommute(OsStruct)

	err := cmd.Run()

	if err != nil {
		LogFunc(fmt.Sprintf("%v,\n", err))
		OsStruct.Exit(1)
	}

}
