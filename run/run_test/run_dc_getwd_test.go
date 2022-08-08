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

// func (osStruct *dcYmlOsGetwdDouble) ReadFile(name string) ([]byte, error) {
// 	osStruct.wasCalled.ReadFile++
// 	return []byte{}, fmt.Errorf(`Failed to Stat() path: '%v'`, name)
// }

func (osStruct *dcYmlOsGetwdDouble) Getwd() (dir string, err error) {
	osStruct.wasCalled.Getwd++
	return `/path/to/dumbclicker`, nil
}
func (*dcYmlOsGetwdDouble) Getenv(string) string { return `` }

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.wasCalled.Stat++
	osStruct.wasCalledWith.Stat = append(osStruct.wasCalledWith.Stat, name)
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

// type execOsDryrunDouble struct {
// 	execOsFailingDouble
// }

// func (execOsStruct *execOsDryrunDouble) Getenv(name string) string {
// 	if name == `TMUX_COMPOSE_DRY_RUN` {
// 		return `1` // dry run
// 	}
// 	return ``
// }

type execOsStructFailingChdir struct {
	execOsFailingDouble
	wasCalled struct {
		Chdir int
	}
}

func (osStruct *execOsStructFailingChdir) Chdir(dir string) (err error) {
	osStruct.wasCalled.Chdir++
	if dir == `/path/to/dumbclicker` {
		err = fmt.Errorf(`Changing directory to: '/path/to/dumbclicker' error: 'not found'`)
	} else {
		err = fmt.Errorf(`Changing directory to: '%v' error: 'wrong directory'`, dir)
	}
	return
}

func TestGetwd(t *testing.T) {

	dcYmlOsStruct := &dcYmlOsGetwdDouble{}
	execOsStruct := &execOsStructFailingChdir{}
	osStruct := &osDouble{}

	runner := run.Runner{
		CmdNameArgs: cmd_name_args.Construct(&cnaOsFailingDouble{}),
		DcYml:       dc_yml.Construct(dcYmlOsStruct),
		Exec:        exec.Construct(execOsStruct),
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
	if execOsStruct.wasCalled.Chdir != 1 {
		t.Errorf(`'execOsStruct.Chdir()' was called not '1' times but: '%v'`, execOsStruct.wasCalled.Chdir)
	}
	stderrExpected := "Changing directory to: '/path/to/dumbclicker' error: 'not found'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing ExecOsStruct.Chdir() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
