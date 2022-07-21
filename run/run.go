package run

import (
	"tmux_compose/dc_config"
)

type LogFuncType func(v ...any)
type CmdNameArgsType func(dcConfigReader dc_config.Reader) (string, []string)

type Runner struct {
	CmdNameArgs    CmdNameArgsType
	DcConfigReader dc_config.Reader
	ExecStruct     execInterface
	OsStruct       *OsStruct
	LogFunc        LogFuncType
}

func (runner Runner) Run() {
	DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	CmdNameArgs := runner.CmdNameArgs
	cmdName, args := CmdNameArgs(DcConfigReader)

	cmd := ExecStruct.GetCommand(cmdName, args...)
	cmd.StdCommute(OsStruct)

	err := cmd.Obj.Run()
	if err != nil {
		LogFunc(err)
	}

}
