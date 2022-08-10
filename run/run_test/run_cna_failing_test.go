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

	_ "embed"
)

type dcYmlOsGetwdRootDouble struct {
	dcYmlOsFailingDouble
}

//go:embed testdata/bash-new-window.gson
var cnaDefaultTemplateContents []byte

//go:embed testdata/dumbclicker/docker-compose.yml
var dcYmlDefault []byte

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) ReadFile(name string) ([]byte, error) {
	if name == `/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Getwd() (dir string, err error) {
	return `/`, nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Stat(name string) (dfi types.FileInfoStruct, err error) {
	if name == `/` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/docker-compose.yml` {
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

func TestRunCnaWorkdirRoot(t *testing.T) {

	dcYml := dc_yml.Construct(&dcYmlOsGetwdRootDouble{})
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

	if os.ExitData.code != 1 {
		t.Errorf(`Workdir '/' provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Workdir '/' called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	stderrExpected := "Get command name and args error: error finding base dir name '/' same length for work dir: '/'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type cnaOsFailingStat struct {
	cnaOsFailingDouble
}

// Stat implements types.DcYmlOsInterface
func (osStruct *cnaOsFailingStat) Stat(name string) (dfi types.FileInfoStruct, err error) {
	return dfi, fmt.Errorf("not found")
}

type configCnaTemplate struct {
	configFailingDouble
}

func (*configCnaTemplate) GetCnaTemplateFname() string {
	return `/path/to/dumbclicker/tmux-compose.json`
}

func TestRunCnaOsFailingStat(t *testing.T) {

	dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{})
	cna := cmd_name_args.Construct(&cnaOsFailingStat{}, &configCnaTemplate{})
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

	if os.ExitData.code != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get command name and args error: error reading file '/path/to/dumbclicker/tmux-compose.json': not found\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type cnaOsFailingReadFile struct {
	ReadFileData struct {
		WasCalled int
		Args      []string
	}
}

// ReadFile implements types.CnaOsInterface
func (cnaOsStruct cnaOsFailingReadFile) ReadFile(name string) ([]byte, error) {
	cnaOsStruct.ReadFileData.WasCalled++
	cnaOsStruct.ReadFileData.Args = append(cnaOsStruct.ReadFileData.Args, name)
	if name == `/path/to/dumbclicker/tmux-compose.json` {
		return []byte{}, fmt.Errorf("permission denied")
	} else {
		return []byte{}, fmt.Errorf(`Wrong file name: '%v' not '/path/to/dumbclicker/tmux-compose.json'`, name)
	}
}

// Stat implements types.DcYmlOsInterface
func (osStruct *cnaOsFailingReadFile) Stat(name string) (dfi types.FileInfoStruct, err error) {
	if name == `/path/to/dumbclicker/tmux-compose.json` {
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

func TestRunCnaOsFailingReadFile(t *testing.T) {

	dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{})
	cnaOsStruct := &cnaOsFailingReadFile{}
	cna := cmd_name_args.Construct(cnaOsStruct, &configCnaTemplate{})
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

	if cnaOsStruct.ReadFileData.WasCalled != 1 {
		if os.ExitData.code != 1 {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called not '1' time(s)  but: '%v'`, os.ExitData.code)
		}
	} else {
		if cnaOsStruct.ReadFileData.Args[0] != `/path/to/dumbclicker/tmux-compose.json` {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called with '%v' arg(s) not with: '%v'`, cnaOsStruct.ReadFileData.Args[0], `/path/to/dumbclicker/tmux-compose.json`)
		}
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get command name and args error: error reading file '/path/to/dumbclicker/tmux-compose.json': permission denied\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

func TestRunCnaOsFailingTmplExecute(t *testing.T) {

	dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{})
	cnaOsStruct := &cnaOsFailingReadFile{}
	cna := cmd_name_args.Construct(cnaOsStruct, &configCnaTemplate{})
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

	if cnaOsStruct.ReadFileData.WasCalled != 1 {
		if os.ExitData.code != 1 {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called not '1' time(s)  but: '%v'`, os.ExitData.code)
		}
	} else {
		if cnaOsStruct.ReadFileData.Args[0] != `/path/to/dumbclicker/tmux-compose.json` {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called with '1' time(s)  but: '%v'`, cnaOsStruct.ReadFileData.Args[0])
		}
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing exec template was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing exec template was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing exec template made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}

	stderrExpected := fmt.Sprintf("Get command name and args error: error executing template '%s' on with vars from: '%s': %v\n", cnaDefaultTemplateContents, dcYmlDefault, ``)
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing exec template made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}
