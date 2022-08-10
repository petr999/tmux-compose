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

type dcYmlOsGetwdRootDouble struct {
	dcYmlOsFailingDouble
	GetwdData struct{ WascalledTimes int }
}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) ReadFile(name string) ([]byte, error) {
	if name == `/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Getwd() (dir string, err error) {
	osStruct.GetwdData.WascalledTimes++
	return `/`, nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	if name == `/` {
		return types.DcFileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/docker-compose.yml` {
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

func TestRunCnaWorkdirRoot(t *testing.T) {

	dcYmlOsStruct := &dcYmlOsGetwdRootDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{}, &configFailingDouble{})
	execOsStruct := &execOsFailingDouble{}
	exec := exec.Construct(execOsStruct)

	os := &osDouble{}
	logger := logger.Construct(execOsStruct.GetStdHandles())

	runner := run.Runner{
		CmdNameArgs: cna,
		DcYml:       dcYml,
		Exec:        exec,
		Os:          os,
		Logger:      logger,
	}

	runner.Run()

	if dcYmlOsStruct.GetwdData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called not '1' time but: '%v'`, dcYmlOsStruct.GetwdData.WascalledTimes)
	}
	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get command name and args error: error finding base dir name '/' same length for work dir: '/'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
