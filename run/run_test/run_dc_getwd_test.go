package run_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
	"tmux_compose/types"
)

type dcYmlOsGetwdDouble struct {
	dcYmlOsReadFileDouble
	wasCalled struct {
		ReadFile int
		Getwd    int
		Stat     int
	}
	wasCalledWith struct {
		Stat []string
	}
}

func (osStruct *dcYmlOsGetwdDouble) Getwd() (dir string, err error) {
	osStruct.wasCalled.Getwd++
	return `/path/to/dumbclicker`, nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdDouble) Stat(name string) (dfi types.FileInfoStruct, err error) {
	osStruct.wasCalled.Stat++
	osStruct.wasCalledWith.Stat = append(osStruct.wasCalledWith.Stat, name)
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

type cnaOsGetpwd struct {
	cnaOsFailingDouble
	StatData struct {
		Args []string
	}
}

func (cnaOsStruct *cnaOsGetpwd) Stat(name string) (types.FileInfoStruct, error) {
	isDir := false
	for _, arg := range cnaOsStruct.StatData.Args {
		if arg == filepath.Base(name) {
			isDir = true
			continue
		}
	}
	return types.FileInfoStruct{
		IsDir: func() bool {
			return isDir
		},
		IsFile: func() bool {
			return !isDir
		},
	}, nil
}

func TestGetwd(t *testing.T) {

	dcYmlOsStruct := &dcYmlOsGetwdDouble{}
	execOsStruct := &execOsStructFailingChdir{}
	osStruct := &osDouble{}
	configStruct := config.Construct(ConfigOsDouble{})

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsGetpwd{}, configStruct),
		DcYml:       dc_yml.Construct(dcYmlOsStruct, configStruct),
		Exec:        exec.Construct(execOsStruct, configStruct),
		Os:          osStruct,
		Logger:      logger.Construct(execOsStruct.GetStdHandles()),
	}

	runner.Run()

	if dcYmlOsStruct.wasCalled.Getwd != 1 {
		t.Errorf(`dcYmlOsStruct.Getwd was called '%v' times instead of '1'`, dcYmlOsStruct.wasCalled.Getwd)
	}
	if osStruct.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, osStruct.ExitData.code)
	}
	if osStruct.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, osStruct.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Changing directory to: '/path/to/dumbclicker' error: 'not found'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing ExecOsStruct.Chdir() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
