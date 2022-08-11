package run_test

import (
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
func (osStruct *dcYmlOsFailingParse) Stat(name string) (dfi types.FileInfoStruct, err error) {
	if name == `/path/to/dumbclicker` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/path/to/dumbclicker/docker-compose.yml` {
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

func (dcYml *dcYmlOsFailingParse) Getwd() (string, error) {
	return `/path/to/dumbclicker`, nil
}

func (dcYml *dcYmlOsFailingParse) SetDcContents(dcContents string) {
	dcYml.dcContents = dcContents
}

func TestRunDcParse(t *testing.T) {
	for _, dcStderrs := range []map[string]string{
		{``: `Get docker-compose config error: no service names in config: '/path/to/dumbclicker/docker-compose.yml'`},
		{`v: [A,`: `Get docker-compose config error: get services from '/path/to/dumbclicker/docker-compose.yml': yaml: line 1: did not find expected node content`},
		{"version: \"3.7\"\n\nservices:\n": `Get docker-compose config error: get services from '/path/to/dumbclicker/docker-compose.yml': services are not a list of strings: '[{version 3.7} {services <nil>}]'`},
		{"version: \"3.7\"\n\nservices:\n  \"\":\n      image: nginx:latest": `Get docker-compose config error: empty or inappropriate service name '' in config: '/path/to/dumbclicker/docker-compose.yml'`},
		{"version: \"3.7\"\n\nservices:\n    1:\n      image: nginx:latest": `Get docker-compose config error: get services from '/path/to/dumbclicker/docker-compose.yml': service name is not a string: '1'`},
		{"version: \"3.7\"\n\nservices:\n    \"-\":\n      image: nginx:latest": `Get docker-compose config error: empty or inappropriate service name '-' in config: '/path/to/dumbclicker/docker-compose.yml'`},
	} {
		for dcContents, stderrExpected := range dcStderrs {

			dcYmlOsStruct := &dcYmlOsFailingParse{}
			dcYmlOsStruct.SetDcContents(dcContents)
			execOsStruct := &execOsFailingDouble{}
			osStruct := &osDouble{}
			stderrExpected = stderrExpected + "\n"
			configStruct := config.Construct(ConfigOsDouble{})

			runner := run.Runner{
				CmdNameArgs: cmd_name_args.Construct(&cnaOsFailingDouble{}, configStruct),
				DcYml:       dc_yml.Construct(dcYmlOsStruct, configStruct),
				Exec:        exec.Construct(execOsStruct, configStruct),
				Os:          osStruct,
				Logger:      logger.Construct(execOsStruct.GetStdHandles()),
			}

			runner.Run()

			if dcYmlOsStruct.wasCalled.ReadFile != 1 {
				t.Errorf(`dcYml.Get was called '%v' times instead of '1'`, dcYmlOsStruct.wasCalled.ReadFile)
			}
			if osStruct.ExitData.code != 1 {
				t.Errorf(`Failing dcYml.Get was provided not '1' to Runner.Os.Exit exit code but: '%v'`, osStruct.ExitData.code)
			}
			if osStruct.ExitData.wasCalledTimes != 1 {
				t.Errorf(`Failing dcYml.Get was called Runner.Os.Exit not '1' time: '%v'`, osStruct.ExitData.code)
			}
			if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
				t.Errorf(`Failing dcYml.Get made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
			}
			if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
				t.Errorf(`Failing dcYml.Get parse '%v' made stderr '%s' not equal to: '%s'`, dcContents, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
			}
		}
	}

}
