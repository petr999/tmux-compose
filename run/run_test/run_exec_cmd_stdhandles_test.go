package run_test

import (
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
	"tmux_compose/types"
)

type execOsStructChdir struct {
	execOsFailingDouble
}

func (osStruct *execOsStructChdir) Chdir(dir string) (err error) { return }

func TestExecCmdStdhandles(t *testing.T) {

	execOsStruct := &execOsStructChdir{}
	osStruct := &osDouble{}
	exec := exec.Construct(execOsStruct)

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsFailingDouble{}),
		DcYml:       dc_yml.Construct(&dcYmlOsGetwdDouble{}),
		Exec:        exec,
		Os:          osStruct,
		Logger:      logger.Construct(execOsStruct.GetStdHandles()),
	}

	runner.Run()

	if exec.GetCommand(types.CmdNameArgsValueType{}).Stdhandles.Stdout != execOsStruct.GetStdHandles().Stdout {
		t.Errorf(`stdout of command to run: '%p' was not assigned to stdout of os: '%p'`, exec.GetCommand(types.CmdNameArgsValueType{}).Stdhandles.Stdout, execOsStruct.GetStdHandles().Stdout)
	}

	if exec.GetCommand(types.CmdNameArgsValueType{}).Stdhandles.Stderr != execOsStruct.GetStdHandles().Stderr {
		t.Errorf(`stderr of command to run: '%p' was not assigned to stderr of os: '%p'`, exec.GetCommand(types.CmdNameArgsValueType{}).Stdhandles.Stderr, execOsStruct.GetStdHandles().Stderr)
	}
}
