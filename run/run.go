package run

import (
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
)

const DryRunEnvVarName = `TMUX_COMPOSE_DRY_RUN`

type LogFuncType func(v ...any)
type CmdNameArgsType func(dcConfigReader dc_config.Reader) (string, []string)

type Runner struct {
	CmdNameArgs    CmdNameArgsType
	DcConfigReader dc_config.Reader
	ExecStruct     exec.ExecInterface
	OsStruct       *exec.OsStruct
	LogFunc        LogFuncType
}

func (runner *Runner) Run() {
	DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	CmdNameArgs := runner.CmdNameArgs
	cmdName, args := CmdNameArgs(DcConfigReader)

	ExecStruct.MakeCommand(cmdName, args...)
	cmd := ExecStruct.GetCommand()
	cmd.StdCommute(OsStruct)

	err := cmd.Run()

	if err != nil {
		LogFunc(err)
		OsStruct.Exit(1)
	}

}