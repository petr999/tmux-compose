package run_test

import (
	_ "embed"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
)

//go:embed testdata/sample.sh
var dryRunOutput []byte

type configOsDryRun struct{}

func (config *configOsDryRun) Getenv(name string) string {
	if name == `TMUX_COMPOSE_DRY_RUN` {
		return `1`
	}
	return ``
}

func TestExecDryRun(t *testing.T) {

	configOsStruct := &configOsDryRun{}
	configStruct := config.Construct(configOsStruct)
	execOsStruct := &execOsStructChdir{}
	os := &osDouble{}

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsGetpwd{}, configStruct),
		DcYml:       dc_yml.Construct(&dcYmlOsGetwdDouble{}, configStruct),
		Exec:        exec.Construct(execOsStruct, configStruct),
		Os:          os,
		Logger:      logger.Construct(execOsStruct.GetStdHandles()),
	}

	runner.Run()

	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Dry run called Runner.Os.Exit not '1' time(s): '%v'`, os.ExitData.wasCalledTimes)
	}
	if os.ExitData.code != 0 {
		t.Errorf(`Dry run provided not '0' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stderr.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr not empty: '%s'`, execOsStruct.StdHandlesDouble.Stderr)
	}
	stdoutExpected := string(dryRunOutput) + "\n"
	if execOsStruct.StdHandlesDouble.Stdout.String() != stdoutExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stdout, stdoutExpected)
	}

}

func TestExecFail(t *testing.T) {}
