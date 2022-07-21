package main

import (
	"log"
	"os"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_config"
	"tmux_compose/run"
	"tmux_compose/runexec"
)

var runner run.Runner

func init() {
	runner = run.Runner{
		CmdNameArgs:    cmd_name_args.CmdNameArgs,
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     runexec.ExecStruct{},
		OsStruct:       &runexec.OsStruct{Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin},
		LogFunc:        log.Fatal,
	}
}

func main() {
	runner.Run()
}
