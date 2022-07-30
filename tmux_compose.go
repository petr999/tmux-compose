package main

import (
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_config"
	"tmux_compose/run"
	"tmux_compose/run/exec"
	"tmux_compose/types"
)

var runner run.Runner

func init() {

	execOsStruct := types.MakeOsStruct()
	dcOsStruct := types.MakeOsStruct()

	runner = run.Runner{
		CmdNameArgs:    cmd_name_args.CmdNameArgs,
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     &exec.ExecStruct{},
		OsStruct:       execOsStruct,
		LogFunc:        exec.GetLogFunc(dcOsStruct.Stderr),
	}
}

func main() {
	runner.Run()
}
