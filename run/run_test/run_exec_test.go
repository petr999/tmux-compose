package run_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
	"tmux_compose/types"
)

//go:embed testdata/sample.sh
var dryRunOutput []byte

//go:embed testdata/sample-grid.sh
var dryRunOutputGrid []byte

//go:embed testdata/dumbclicker-grid/tmux-compose-template.gson
var cnaTemplateGrid []byte

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
		t.Errorf(`Dry run made stderr not empty: '%s'`, execOsStruct.StdHandlesDouble.Stderr)
	}
	stdoutExpected := string(dryRunOutput) + "\n"
	if execOsStruct.StdHandlesDouble.Stdout.String() != stdoutExpected {
		t.Errorf(`Dry run made stdout '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stdout, stdoutExpected)
	}

}

type configOsDryRunGrid struct {
	ConfigOsCnaOsGetenvFile
	configOsDryRun
}

func (config *configOsDryRunGrid) Getenv(name string) (val string) {
	val = config.configOsDryRun.Getenv(name)
	if len(val) == 0 {
		val = config.ConfigOsCnaOsGetenvFile.Getenv(name)
	}
	return
}

type cnaOsGrid struct {
	cnaOsStatFile
}

func (cna *cnaOsGrid) ReadFile(name string) ([]byte, error) {
	return cnaTemplateGrid, nil
}

type dcYmlOsGrid struct {
	dcYmlOsGetwdDouble
}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGrid) ReadFile(name string) ([]byte, error) {
	if name == `/path/to/dumbclicker-grid/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

func (osStruct *dcYmlOsGrid) Getwd() (dir string, err error) {
	osStruct.wasCalled.Getwd++
	return `/path/to/dumbclicker-grid`, nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGrid) Stat(name string) (dfi types.FileInfoStruct, err error) {
	osStruct.wasCalled.Stat++
	osStruct.wasCalledWith.Stat = append(osStruct.wasCalledWith.Stat, name)
	if name == `/path/to/dumbclicker-grid` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/path/to/dumbclicker-grid/docker-compose.yml` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return false
			},
			IsFile: func() bool {
				return true
			},
		}, nil
	}
	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
}

func TestExecDryRunGrid(t *testing.T) {

	configOsStruct := &configOsDryRunGrid{}
	configStruct := config.Construct(configOsStruct)
	execOsStruct := &execOsStructChdir{}
	os := &osDouble{}

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsGrid{}, configStruct),
		DcYml:       dc_yml.Construct(&dcYmlOsGrid{}, configStruct),
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
		t.Errorf(`Dry run made stderr not empty: '%s'`, execOsStruct.StdHandlesDouble.Stderr)
	}
	stdoutExpected := string(dryRunOutputGrid) + "\n"
	if execOsStruct.StdHandlesDouble.Stdout.String() != stdoutExpected {
		t.Errorf(`Dry run made stdout '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stdout, stdoutExpected)
	}

}

func TestExecFail(t *testing.T) {

	configOsStruct := &ConfigOsDouble{}
	configStruct := config.Construct(configOsStruct)
	execOsStruct := &execOsStructChdir{}
	execStruct := exec.Construct(execOsStruct, configStruct)
	os := &osDouble{}

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsGetpwd{}, configStruct),
		DcYml:       dc_yml.Construct(&dcYmlOsGetwdDouble{}, configStruct),
		Exec:        execStruct,
		Os:          os,
		Logger:      logger.Construct(execOsStruct.GetStdHandles()),
	}

	runner.Run()

	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Run called Runner.Os.Exit not '1' time(s): '%v'`, os.ExitData.wasCalledTimes)
	}
	if os.ExitData.code != 1 {
		t.Errorf(`Run provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if execStruct.Cmd.RunData.WasCalledTimes != 1 {
		t.Errorf(`Run was called  not '1' time(s): '%v'`, execStruct.Cmd.RunData.WasCalledTimes)
	}
	act, _ := json.Marshal(execStruct.Cmd.RunData.Cnas)
	actStr, expectedStr := string(act), `[{"Shebang":"bash","Workdir":"/path/to/dumbclicker","Cmd":"tmux","Args":["new","-s","dumbclicker-compose","\n  docker-compose up\n  bash -l\n",";","neww","-n","dumbclicker_nginx_1","\n  PID=0\n  try_next=1\n  trap '\n    echo \"trap pid: ${PID}\"\n    kill -INT $PID\n    try_next=\"\"\n  ' SIGINT\n  while [ 'x1' == \"x${try_next}\" ]; do\n    bash -lc '\n      docker attach dumbclicker_nginx_1\n      sleep 1\n    ' \u0026\n    PID=$!\n    echo \"pid: ${PID}\"\n    wait $PID\n  done\n  trap - SIGINT\n  bash -l\n",";","neww","-n","dumbclicker_h2o_1","\n  PID=0\n  try_next=1\n  trap '\n    echo \"trap pid: ${PID}\"\n    kill -INT $PID\n    try_next=\"\"\n  ' SIGINT\n  while [ 'x1' == \"x${try_next}\" ]; do\n    bash -lc '\n      docker attach dumbclicker_h2o_1\n      sleep 1\n    ' \u0026\n    PID=$!\n    echo \"pid: ${PID}\"\n    wait $PID\n  done\n  trap - SIGINT\n  bash -l\n",";","neww","-n","dumbclicker_dumbclicker_1","\n  PID=0\n  try_next=1\n  trap '\n    echo \"trap pid: ${PID}\"\n    kill -INT $PID\n    try_next=\"\"\n  ' SIGINT\n  while [ 'x1' == \"x${try_next}\" ]; do\n    bash -lc '\n      docker attach dumbclicker_dumbclicker_1\n      sleep 1\n    ' \u0026\n    PID=$!\n    echo \"pid: ${PID}\"\n    wait $PID\n  done\n  trap - SIGINT\n  bash -l\n"]}]`
	if actStr != expectedStr {
		t.Errorf(`Run was called with '%s' not with '%s' arg`, actStr, expectedStr)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Run made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "open terminal failed: not a terminal\nexit status 1" + "\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Run made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}

}
