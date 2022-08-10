package run_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"
	"tmux_compose/types"

	"tmux_compose/run"

	_ "embed"
)

//go:embed testdata/dumbclicker/docker-compose.yml
var dcContents []byte

type dcYmlOsReadFileDouble struct{}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsReadFileDouble) ReadFile(name string) ([]byte, error) {
	if name == `/path/to/dumbclicker/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

type dcYmlOsGetenvDouble struct {
	dcYmlOsFailingDouble
	dcYmlOsReadFileDouble
	wasCalled struct {
		ReadFile int
		Getwd    int
		Stat     int
	}
}

func (osStruct *dcYmlOsGetenvDouble) ReadFile(name string) ([]byte, error) {
	osStruct.wasCalled.ReadFile++
	return osStruct.dcYmlOsReadFileDouble.ReadFile(name)
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvDouble) Getwd() (dir string, err error) {
	osStruct.wasCalled.Getwd++
	err = fmt.Errorf("unimplemented")
	return
}

func (osStruct *dcYmlOsGetenvDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.wasCalled.Stat++
	return dfi, fmt.Errorf("unimplemented")
}

type dcYmlOsGetenvToFileDouble struct {
	dcYmlOsGetenvDouble
	GetenvData struct{ WascalledTimes int }
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetenvToFileDouble) Getenv(key string) string {
	osStruct.GetenvData.WascalledTimes++
	if key == `TMUX_COMPOSE_DC_YML` {
		return `/path/to/dumbclicker/docker-compose.yml`
	}
	return ``
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Fqfn-resolving DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Stat != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: unimplemented\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}

}

func TestRunDcOsGetenvDir(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsGetenvToDirDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Stat != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
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
	stderrExpected := "Get docker-compose config error: unimplemented\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

// type dcYmlOsFailingChdirDouble struct {
// 	dcYmlOsGetenvToFileDouble
// 	ChdirData struct{ Dir string }
// }

// // Stat implements types.DcYmlOsInterface
// func (osStruct *dcYmlOsFailingChdirDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
// 	osStruct.wasCalled.Stat++
// 	osStruct.StatData.Names = append(osStruct.StatData.Names, name)
// 	if name == `/path/to/dumbclicker` {
// 		return types.DcFileInfoStruct{
// 			IsDir: func() bool {
// 				return false
// 			},
// 			IsFile: func() bool {
// 				return true
// 			},
// 		}, nil
// 	}
// 	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
// }

// func TestRunDcFailingChdir(t *testing.T) { // AndStdHandles {
// 	// tle := getTestLogfuncExitType()
// 	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

// 	dcYmlOsStruct := &dcYmlOsFailingChdirDouble{}
// 	dcYml := dc_yml.Construct(dcYmlOsStruct)
// 	cna := cmd_name_args.Construct(&cnaOsFailingDouble{})
// 	execOsStruct := &execOsFailingDouble{}
// 	exec := exec.Construct(execOsStruct)

// 	os := &osDouble{}
// 	logger := logger.Construct(execOsStruct.GetStdHandles())

// 	runner := run.Runner{
// 		CmdNameArgs: cna,
// 		DcYml:       dcYml,
// 		Exec:        exec,
// 		Os:          os,
// 		Logger:      logger,
// 	}

// 	runner.Run()

// 	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
// 		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
// 	}
// 	if dcYmlOsStruct.wasCalled.Stat != 1 {
// 		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
// 	}
// 	if dcYmlOsStruct.wasCalled.Chdir != 1 {
// 		t.Errorf(`Failing DcOsStruct.Chdir() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.Chdir)
// 	}
// 	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
// 		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
// 	}
// 	if dcYmlOsStruct.wasCalled.Getwd != 0 {
// 		t.Errorf(`Failing DcOsStruct.Getwd() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
// 	}

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
// }

type dcYmlOsFailingReadFileDouble struct {
	dcYmlOsGetenvToFileDouble
	ReadFileData struct{ Name string }
}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingReadFileDouble) ReadFile(name string) ([]byte, error) {
	osStruct.wasCalled.ReadFile++
	osStruct.ReadFileData.Name = name
	return []byte{}, fmt.Errorf("Failed to ReadFile() from: '%v'", name)
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingReadFileDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.wasCalled.Stat++ // .StatData.wasCalled++
	if name == `/path/to/dumbclicker/docker-compose.yml` {
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

func TestRunDcFailingReadFile(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingReadFileDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.Stat != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' times but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 1 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '1' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
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
	stderrExpected := "Get docker-compose config error: get services from '/path/to/dumbclicker/docker-compose.yml': Failed to ReadFile() from: '/path/to/dumbclicker/docker-compose.yml'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type dcYmlOsFailingStatInputDouble struct {
	StatData struct {
		Names     []string
		wasCalled int
	}
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingStatInputDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.StatData.wasCalled++
	osStruct.StatData.Names = append(osStruct.StatData.Names, name)
	return dfi, fmt.Errorf("Failed to Stat() path: '%v'", name)
}

type dcYmlOsFailingStatInputDirDouble struct {
	dcYmlOsGetenvToDirDouble
	dcYmlOsFailingStatInputDouble
}

func TestRunDcFailingStatInputDir(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingStatInputDirDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.StatData.wasCalled != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' time but: '%v'`, dcYmlOsStruct.StatData.wasCalled)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}
	names, _ := json.Marshal(dcYmlOsStruct.StatData.Names)
	namesExpected, _ := json.Marshal([]string{`/path/to/dumbclicker`})
	if string(names) != string(namesExpected) {
		t.Errorf(`Failing DcOsStruct.Stat() was called not with '%s'  but with: '%s'`, namesExpected, names)
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: Failed to Stat() path: '/path/to/dumbclicker'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type dcYmlOsFailingStatInputFileDouble struct {
	dcYmlOsGetenvToFileDouble
	dcYmlOsFailingStatInputDouble
}

func TestRunDcFailingStatInputFile(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingStatInputFileDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}
	if dcYmlOsStruct.StatData.wasCalled != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' times but: '%v'`, dcYmlOsStruct.StatData.wasCalled)
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: Failed to Stat() path: '/path/to/dumbclicker/docker-compose.yml'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

// DcYml supplied is neither dir nor file
type dcYmlOsStatIsOtherInputDouble struct {
	dcYmlOsGetenvToDirDouble
	StatData struct {
		Names []string
	}
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsStatIsOtherInputDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.wasCalled.Stat++
	osStruct.StatData.Names = append(osStruct.StatData.Names, name)
	if name == `/path/to/dumbclicker` {
		return types.DcFileInfoStruct{
			IsDir: func() bool {
				return false
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	}
	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
}

func TestRunDcStatIsOtherInputFile(t *testing.T) { // AndStdHandles {

	dcYmlOsStruct := &dcYmlOsStatIsOtherInputDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}
	if dcYmlOsStruct.wasCalled.Stat != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '1' times but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: not a dir or file: '/path/to/dumbclicker'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

// Both dir and file are directories
type dcYmlOsStatIsDirBothInputDouble struct {
	dcYmlOsGetenvToDirDouble
	StatData struct{ Names []string }
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsStatIsDirBothInputDouble) Stat(name string) (dfi types.DcFileInfoStruct, err error) {
	osStruct.wasCalled.Stat++ // .StatData.wasCalled++
	osStruct.StatData.Names = append(osStruct.StatData.Names, name)
	if (name == `/path/to/dumbclicker`) || (name == `/path/to/dumbclicker/docker-compose.yml`) {
		return types.DcFileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	}
	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
}

func TestRunDcStatIsDirBothInputFile(t *testing.T) { // AndStdHandles {

	dcYmlOsStruct := &dcYmlOsStatIsDirBothInputDouble{}
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

	if dcYmlOsStruct.GetenvData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getenv() was called not '1' time but: '%v'`, dcYmlOsStruct.GetenvData.WascalledTimes)
	}
	if dcYmlOsStruct.wasCalled.ReadFile != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.ReadFile)
	}
	if dcYmlOsStruct.wasCalled.Getwd != 0 {
		t.Errorf(`Failing DcOsStruct.ReadFile() was called not '0' time but: '%v'`, dcYmlOsStruct.wasCalled.Getwd)
	}
	if dcYmlOsStruct.wasCalled.Stat != 2 {
		t.Errorf(`Failing DcOsStruct.Stat() was called not '2' times but: '%v'`, dcYmlOsStruct.wasCalled.Stat)
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Stat() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Stat() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: not a file: '/path/to/dumbclicker/docker-compose.yml'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Stat() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
