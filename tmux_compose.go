package main

import (
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/os_struct"
	"tmux_compose/run"
)

var runner run.Runner

func init() {

	runner = run.Runner{
		CmdNameArgs: cmd_name_args.Construct(os_struct.CnaOsStruct{}),
		DcYml:       dc_yml.Construct(os_struct.DcYmlOsStruct{}),
		Exec:        exec.Construct(os_struct.ExecOsStruct{}),
		Os:          os_struct.RunnerOsStruct{},
		Logger:      logger.Construct(logger.GetStdHandles()),
	}
}

func main() {
	runner.Run()
}
