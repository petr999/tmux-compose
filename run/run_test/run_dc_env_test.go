package run_test

import (
	"fmt"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"

	"tmux_compose/run"

	_ "embed"
)

//go:embed testdata/dumbclicker/docker-compose.yml
var dcContents []byte

type dcYmlOsGetenvDouble struct {
	dcYmlOsFailingDouble
	wasCalled struct {
		Chdir    int
		ReadFile int
		Getwd    int
	}
}

// Chdir implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvDouble) Chdir(dir string) error {
	osStruct.wasCalled.Chdir++
	if dir == `/path/to/dumbclicker` {
		return nil
	}
	return fmt.Errorf("Wrong dir to Chdir(): '%v'", dir)
}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvDouble) ReadFile(name string) ([]byte, error) {
	osStruct.wasCalled.ReadFile++
	if name == `/path/to/dumbclicker/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvDouble) Getwd() (dir string, err error) {
	osStruct.wasCalled.Getwd++
	err = fmt.Errorf("unimplemented")
	return
}

type dcYmlOsGetenvToFileDouble struct {
	dcYmlOsGetenvDouble
	GetenvData struct{ WascalledTimes int }
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvToFileDouble) Getenv(key string) (val string) {
	osStruct.GetenvData.WascalledTimes++
	val = ``
	if key == `TMUX_COMPOSE_DC_YML` {
		val = `/path/to/dumbclicker/docker-compose.yml`
	}
	return
}

type dcYmlOsGetenvToDirDouble struct {
	dcYmlOsGetenvDouble
	GetenvData struct{ WascalledTimes int }
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvToDirDouble) Getenv(key string) (val string) {
	osStruct.GetenvData.WascalledTimes++
	val = ``
	if key == `TMUX_COMPOSE_DC_YML` {
		val = `/path/to/dumbclicker`
	}
	return
}

func TestRunDcOsGetenvFile(t *testing.T) { // AndStdHandles {

	dcYmlOsStruct := &dcYmlOsGetenvToFileDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{})
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Chdir != 1 {
		t.Errorf(`Failing DcOsStruct.Chdir() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Chdir)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}

	// if os.ExitData.code != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	// }
	// if os.ExitData.wasCalledTimes != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	// }
	// if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	// }
	// stderrExpected := "exec: no command\n"
	// if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	// }

}

func TestRunDcOsGetenvDir(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsGetenvToDirDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{})
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Chdir != 1 {
		t.Errorf(`Failing DcOsStruct.Chdir() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Chdir)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}

	// if os.ExitData.code != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	// }
	// if os.ExitData.wasCalledTimes != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	// }
	// if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	// }
	// stderrExpected := "exec: no command\n"
	// if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	// }
}

type dcYmlOsFailingChdirDouble struct {
	dcYmlOsGetenvToFileDouble
	ChdirData struct{ Dir string }
}

// Chdir implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingChdirDouble) Chdir(dir string) error {
	osStruct.wasCalled.Chdir++
	osStruct.ChdirData.Dir = dir
	return fmt.Errorf("Failed to Chdir() to: '%v'", dir)
}

func TestRunDcFailingChdir(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingChdirDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{})
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Chdir != 1 {
		t.Errorf(`Failing DcOsStruct.Chdir() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Chdir)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}

	// if os.ExitData.code != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	// }
	// if os.ExitData.wasCalledTimes != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	// }
	// if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	// }
	// stderrExpected := "exec: no command\n"
	// if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	// }
}

type dcYmlOsFailingReadFileDouble struct {
	dcYmlOsGetenvToFileDouble
	ReadFileData struct{ Name string }
}

// Chdir implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingReadFileDouble) ReadFile(name string) ([]byte, error) {
	osStruct.wasCalled.ReadFile++
	osStruct.ReadFileData.Name = name
	return []byte{}, fmt.Errorf("Failed to ReadFile() from: '%v'", name)
}

func TestRunDcFailingReadFile(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingReadFileDouble{}
	dcYml := dc_yml.Construct(dcYmlOsStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{})
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Chdir != 1 {
		t.Errorf(`Failing DcOsStruct.Chdir() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Chdir)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}

	// if os.ExitData.code != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	// }
	// if os.ExitData.wasCalledTimes != 1 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	// }
	// if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	// }
	// stderrExpected := "exec: no command\n"
	// if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
	// 	t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	// }
}
