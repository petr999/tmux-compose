package run_test

import (
	"fmt"
	"strings"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/run"
)

type execOsStructFailingChdir struct {
	execOsFailingDouble
	chdirData struct {
		wasCalled int
		dirs      []string
	}
}

func (osStruct *execOsStructFailingChdir) Chdir(dir string) (err error) {
	osStruct.chdirData.wasCalled++
	osStruct.chdirData.dirs = append(osStruct.chdirData.dirs, dir)

	if dir == `/path/to/dumbclicker` {
		err = fmt.Errorf(`Changing directory to: '/path/to/dumbclicker' error: 'not found'`)
	} else {
		err = fmt.Errorf(`Changing directory to: '%v' error: 'wrong directory'`, dir)
	}

	return
}

func TestRunDcFailingChdir(t *testing.T) {
	configOsStruct := &ConfigOsDouble{}
	configStruct := config.Construct(configOsStruct)

	dcYmlOsStruct := &dcYmlOsGetwdDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct, configStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{}, configStruct)
	execOsStruct := &execOsStructFailingChdir{}
	exec := exec.Construct(execOsStruct, configStruct)

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

	if execOsStruct.chdirData.wasCalled != 1 {
		t.Errorf(`Failing execOsStruct.Chdir() was called not '1' time but: '%v'`, execOsStruct.chdirData.wasCalled)
	}
	if strings.Join(execOsStruct.chdirData.dirs, ":") != `/path/to/dumbclicker` {
		t.Errorf(`Failing execOsStruct.Chdir() was called with '%v' arg(s) not with: '%v'`, execOsStruct.chdirData.dirs, `/path/to/dumbclicker`)
	}
	if os.ExitData.code != 1 {
		t.Errorf(`Failing execOsStruct.Chdir() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing execOsStruct.Chdir() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing execOsStruct.Chdir() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "exec: no command"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing execOsStruct.Chdir() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected+"\n")
	}
}
