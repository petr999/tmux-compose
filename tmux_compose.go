package main

import (
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/os_struct"
	"tmux_compose/run"
)

var runner run.Runner

func init() {

	config := config.Construct(os_struct.ConfigOsStruct{})

	runner = run.Runner{
		CmdNameArgs: cmd_name_args.Construct(os_struct.CnaOsStruct{}, config),
		DcYml:       dc_yml.Construct(os_struct.DcYmlOsStruct{}, config),
		Exec:        exec.Construct(os_struct.ExecOsStruct{}, config),
		Os:          &os_struct.RunnerOsStruct{},
		Logger:      logger.Construct(logger.GetStdHandles()),
	}
}

func main() {
	runner.Run()
}
