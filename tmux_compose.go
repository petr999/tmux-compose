package main

import (
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_config"
	"tmux_compose/run"
	"tmux_compose/run/exec"
)

var runner run.Runner

func init() {

	osStruct := exec.MakeOsStruct()

	runner = run.Runner{
		CmdNameArgs:    cmd_name_args.CmdNameArgs,
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     &exec.ExecStruct{},
		OsStruct:       osStruct,
		LogFunc:        exec.GetLogFunc(osStruct.Stderr),
	}
}

func main() {
	runner.Run()
}
