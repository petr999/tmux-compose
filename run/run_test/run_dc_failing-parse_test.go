package run_test

import (
	"fmt"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
	"tmux_compose/types"
)

type dcYmlOsFailingParse struct {
	dcYmlOsFailingDouble
	dcContents string
	wasCalled  struct {
		ReadFile int
	}
}

func (dcYml *dcYmlOsFailingParse) ReadFile(name string) ([]byte, error) {
	dcYml.wasCalled.ReadFile++
	return []byte(dcYml.dcContents), nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingParse) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	if name == `/path/to/dumbclicker` {
		return types.DcFileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/path/to/dumbclicker/docker-compose.yml` {
		return types.DcFileInfoStruct{
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

func (dcYml *dcYmlOsFailingParse) Getwd() (string, error) {
	return `/path/to/dumbclicker`, nil
}

func (dcYml *dcYmlOsFailingParse) SetDcContents(dcContents string) {
	dcYml.dcContents = dcContents
}

func TestRunDcParse(t *testing.T) {

	dcYmlOsStruct := &dcYmlOsFailingParse{}
	dcYmlOsStruct.SetDcContents(`v: [A,`)
	execOsStruct := &execOsFailingDouble{}
	osStruct := &osDouble{}

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsFailingDouble{}),
		DcYml:       dc_yml.Construct(dcYmlOsStruct),
		Exec:        exec.Construct(execOsStruct),
		Os:          osStruct,
		Logger:      logger.Construct(execOsStruct.GetStdHandles()),
	}

	runner.Run()

	if dcYmlOsStruct.wasCalled.ReadFile != 1 {
		t.Errorf(`dcYmlOsStruct.ReadFile was called '%v' times instead of '1'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if osStruct.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, osStruct.ExitData.code)
	}
	if osStruct.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called Runner.Os.Exit not '1' time: '%v'`, osStruct.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: 'yaml: line 1: did not find expected node content'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing dcYmlOsStruct parse made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
