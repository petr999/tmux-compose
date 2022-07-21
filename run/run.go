package run

import (
	"tmux_compose/dc_config"
)

type LogFunc func(v ...any)

type Runner struct {
	DcConfigReader dc_config.Reader
	ExecStruct     execInterface
	OsStruct       *OsStruct
	LogFunc        LogFunc
}

func (runner Runner) Run() {
	DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	cmdName, args := cmdNameArgs(DcConfigReader)

	cmd := ExecStruct.GetCommand(cmdName, args...)
	cmd.StdCommute(OsStruct)

	err := cmd.Obj.Run()
	if err != nil {
		LogFunc(err)
	}

}
