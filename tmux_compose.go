package main

import (
	"os"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_config"
	"tmux_compose/run"
	"tmux_compose/run/exec"
)

var runner run.Runner

func init() {

	runner = run.Runner{
		CmdNameArgs:    cmd_name_args.CmdNameArgs,
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     &exec.ExecStruct{},
		OsStruct:       &exec.OsStruct{Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin, Exit: os.Exit, Getenv: os.Getenv},
		LogFunc:        exec.LogFunc,
	}
}

func main() {
	runner.Run()
}
