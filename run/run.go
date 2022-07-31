package run

import (
	"fmt"
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
	"tmux_compose/types"
)

const DryRunEnvVarName = `TMUX_COMPOSE_DRY_RUN`

type LogFuncType func(s string)

type CmdNameArgsFunc func(dcConfigReader dc_config.ReaderInterface) (types.CmdNameArgsType, error)

type Runner struct {
	CmdNameArgs    CmdNameArgsFunc
	DcConfigReader dc_config.ReaderInterface
	ExecStruct     exec.ExecInterface
	OsStruct       *types.OsStruct
	LogFunc        LogFuncType
}

func (runner *Runner) Run() {
	DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	CmdNameArgs := runner.CmdNameArgs
	cmdNameArgs, err := CmdNameArgs(DcConfigReader)
	if err != nil {
		LogFunc(fmt.Sprintf("%v,\n", err))
		OsStruct.Exit(1)
	}

	ExecStruct.MakeCommand(&exec.MakeCommandDryRunType{DryRun: OsStruct.Getenv(DryRunEnvVarName), OsStruct: OsStruct},
		cmdNameArgs)
	cmd := ExecStruct.GetCommand()
	cmd.StdCommute(OsStruct)

	err = cmd.Run()

	if err != nil {
		LogFunc(fmt.Sprintf("%v,\n", err))
		OsStruct.Exit(1)
	}

}
